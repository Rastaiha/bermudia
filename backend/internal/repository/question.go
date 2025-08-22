package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Rastaiha/bermudia/internal/domain"
)

const (
	questionsSchema = `
CREATE TABLE IF NOT EXISTS questions (
    id VARCHAR(255) PRIMARY KEY,
    text TEXT NOT NULL,
    knowledge_amount INT4 NOT NULL,
    reward_sources TEXT NOT NULL,
    input_type VARCHAR(255) NOT NULL,
    input_accept TEXT NOT NULL
);
`

	answersSchema = `
CREATE TABLE IF NOT EXISTS answers (
    id VARCHAR(255) PRIMARY KEY,
    user_id INT4 NOT NULL,
    question_id VARCHAR(255) NOT NULL,
    status INT4 NOT NULL,
    file_id VARCHAR(255) NOT NULL,
    filename VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
`
)

type sqlQuestionRepository struct {
	db *sql.DB
}

func NewSqlQuestionRepository(db *sql.DB) (domain.QuestionStore, error) {
	_, err := db.Exec(questionsSchema)
	if err != nil {
		return nil, fmt.Errorf("create questions table: %w", err)
	}
	_, err = db.Exec(answersSchema)
	if err != nil {
		return nil, fmt.Errorf("create answers table: %w", err)
	}
	return sqlQuestionRepository{
		db: db,
	}, nil
}

func (s sqlQuestionRepository) SetQuestion(ctx context.Context, question domain.Question) error {
	rewardSources, err := json.Marshal(question.RewardSources)
	if err != nil {
		return err
	}
	accept, err := json.Marshal(question.InputAccept)
	if err != nil {
		return err
	}

	_, err = s.db.ExecContext(ctx,
		`INSERT INTO questions (id, text, knowledge_amount, reward_sources, input_type, input_accept) 
		 VALUES ($1, $2, $3, $4, $5, $6)
		 ON CONFLICT (id) DO UPDATE SET 
		 text = EXCLUDED.text,
		 knowledge_amount = EXCLUDED.knowledge_amount,
		 reward_sources = EXCLUDED.reward_sources,
		 input_type = EXCLUDED.input_type,
		 input_accept = EXCLUDED.input_accept`,
		n(question.ID), n(question.Text), n(question.KnowledgeAmount), rewardSources, n(question.InputType), accept,
	)
	return err
}

func (s sqlQuestionRepository) GetQuestion(ctx context.Context, id string) (domain.Question, error) {
	var q domain.Question
	var rewardSources, accept []byte
	err := s.db.QueryRowContext(ctx,
		`SELECT id, text, knowledge_amount, reward_sources, input_type, input_accept FROM questions WHERE id = $1`,
		id,
	).Scan(&q.ID, &q.Text, &q.KnowledgeAmount, &rewardSources, &q.InputType, &accept)

	if err != nil {
		return domain.Question{}, fmt.Errorf("failed to get question from db: %w", err)
	}

	if err := json.Unmarshal(rewardSources, &q.RewardSources); err != nil {
		return domain.Question{}, err
	}
	if err := json.Unmarshal(accept, &q.InputAccept); err != nil {
		return domain.Question{}, err
	}

	return q, nil
}

func (s sqlQuestionRepository) GetOrCreateAnswer(ctx context.Context, userId int32, answerID string, questionID string) (domain.Answer, error) {
	var answer domain.Answer
	now := time.Now().UTC()

	err := s.db.QueryRowContext(ctx,
		`INSERT INTO answers (id, user_id, question_id, status, file_id, filename, created_at, updated_at) 
		 VALUES ($1, $2, $3, $4, '', '', $5, $5)
		 ON CONFLICT (id) DO UPDATE SET id = EXCLUDED.id
		 RETURNING id, user_id, question_id, status, file_id, filename, created_at, updated_at`,
		n(answerID), n(userId), n(questionID), domain.AnswerStatusEmpty, now,
	).Scan(&answer.ID, &answer.UserID, &answer.QuestionID, &answer.Status, &answer.FileID, &answer.Filename, &answer.CreatedAt, &answer.UpdatedAt)

	if err != nil {
		return domain.Answer{}, fmt.Errorf("failed to get or create answer: %w", err)
	}

	return answer, nil
}

func (s sqlQuestionRepository) SubmitAnswer(ctx context.Context, answerId string, userId int32, fileID, filename string) (domain.Answer, error) {
	var answer domain.Answer
	now := time.Now().UTC()

	err := s.db.QueryRowContext(ctx,
		`UPDATE answers 
		 SET status = CASE WHEN status = $1 OR status = $2 THEN $3 ELSE status END,
		     file_id = CASE WHEN status = $1 OR status = $2 THEN $4 ELSE file_id END,
		     filename = CASE WHEN status = $1 OR status = $2 THEN $5 ELSE filename END,
		     updated_at = CASE WHEN status = $1 OR status = $2 THEN $6 ELSE updated_at END
		 WHERE id = $7 and user_id = $8
		 RETURNING id, user_id, question_id, status, file_id, filename, created_at, updated_at`,
		domain.AnswerStatusEmpty, domain.AnswerStatusWrong, domain.AnswerStatusPending, fileID, filename, now, answerId, userId,
	).Scan(&answer.ID, &answer.UserID, &answer.QuestionID, &answer.Status, &answer.FileID, &answer.Filename, &answer.CreatedAt, &answer.UpdatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return domain.Answer{}, domain.ErrSubmitToNonExistingAnswer
	}
	if err != nil {
		return domain.Answer{}, fmt.Errorf("db: failed to submit answer: %w", err)
	}

	if answer.Status == domain.AnswerStatusCorrect {
		return domain.Answer{}, domain.ErrSubmitToCorrectAnswer
	}
	if answer.Status == domain.AnswerStatusPending && (answer.FileID != fileID || !answer.UpdatedAt.Equal(now)) {
		return domain.Answer{}, domain.ErrSubmitToPendingAnswer
	}

	return answer, nil
}
