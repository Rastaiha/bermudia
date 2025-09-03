package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/Rastaiha/bermudia/internal/domain"
)

const (
	questionsSchema = `
CREATE TABLE IF NOT EXISTS questions (
    question_id VARCHAR(255) PRIMARY KEY,
    book_id VARCHAR(255) NOT NULL,
    knowledge_amount INT4 NOT NULL,
    reward_source VARCHAR(255)
);
CREATE INDEX IF NOT EXISTS idx_questions_book_id ON questions (book_id);
`
	answersSchema = `
CREATE TABLE IF NOT EXISTS answers (
    user_id INT4 NOT NULL,
    question_id VARCHAR(255) NOT NULL,
    status INT4 NOT NULL,
    file_id VARCHAR(255),
    filename VARCHAR(255),
    text_content TEXT,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    PRIMARY KEY (question_id, user_id)
);
CREATE INDEX IF NOT EXISTS idx_answers_user_question ON answers (user_id, question_id, status);
`
	correctionsSchema = `
CREATE TABLE IF NOT EXISTS corrections (
    id VARCHAR(255) PRIMARY KEY,
    user_id INT4 NOT NULL,
    question_id VARCHAR(255) NOT NULL,
    is_correct BOOLEAN NOT NULL,
    applied BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL
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
	_, err = db.Exec(correctionsSchema)
	if err != nil {
		return nil, fmt.Errorf("create corrections table: %w", err)
	}
	return sqlQuestionRepository{
		db: db,
	}, nil
}

func (s sqlQuestionRepository) BindQuestionsToBook(ctx context.Context, bookId string, questions []domain.BookQuestion) (err error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("start transaction: %w", err)
	}
	defer func() {
		if err != nil {
			err2 := tx.Rollback()
			err = errors.Join(err, err2)
		}
	}()
	_, err = tx.ExecContext(ctx, `DELETE FROM questions WHERE book_id = $1`, bookId)
	if err != nil {
		return fmt.Errorf("delete questions: %w", err)
	}
	for _, q := range questions {
		_, err = tx.ExecContext(ctx, `INSERT INTO questions (question_id, book_id, knowledge_amount, reward_source) VALUES ($1, $2, $3, $4) ;`,
			n(q.QuestionID), n(bookId), q.KnowledgeAmount, n(q.RewardSource),
		)
		if err != nil {
			return fmt.Errorf("insert questions: %w", err)
		}
	}
	return tx.Commit()
}

func (s sqlQuestionRepository) answerColumnsToSelect() string {
	return `user_id, question_id, status, file_id, filename, text_content, created_at, updated_at`
}

func (s sqlQuestionRepository) scanAnswer(row scannable, answer *domain.Answer) error {
	return row.Scan(&answer.UserID, &answer.QuestionID, &answer.Status, &answer.FileID, &answer.Filename, &answer.TextContent, &answer.CreatedAt, &answer.UpdatedAt)
}

func (s sqlQuestionRepository) GetOrCreateAnswer(ctx context.Context, userId int32, questionID string) (domain.Answer, error) {
	var answer domain.Answer
	now := time.Now().UTC()

	err := s.scanAnswer(s.db.QueryRowContext(ctx,
		`INSERT INTO answers (user_id, question_id, status, created_at, updated_at) 
		 VALUES ($1, $2, $3, $4, $4)
		 ON CONFLICT DO UPDATE SET user_id = EXCLUDED.user_id
		 RETURNING `+s.answerColumnsToSelect(),
		n(userId), n(questionID), domain.AnswerStatusEmpty, now, now,
	), &answer)

	if err != nil {
		return domain.Answer{}, fmt.Errorf("failed to get or create answer: %w", err)
	}

	return answer, nil
}

func (s sqlQuestionRepository) SubmitAnswer(ctx context.Context, userId int32, questionId string, fileID, filename, textContent string) (domain.Answer, error) {
	var answer domain.Answer
	now := time.Now().UTC()

	err := s.scanAnswer(s.db.QueryRowContext(ctx,
		`UPDATE answers 
		 SET status = CASE WHEN status = $1 OR status = $2 THEN $3 ELSE status END,
		     file_id = CASE WHEN status = $1 OR status = $2 THEN $4 ELSE file_id END,
		     filename = CASE WHEN status = $1 OR status = $2 THEN $5 ELSE filename END,
		     text_content = CASE WHEN status = $1 OR status = $2 THEN $6 ELSE text_content END,
		     updated_at = CASE WHEN status = $1 OR status = $2 THEN $7 ELSE updated_at END
		 WHERE user_id = $8 AND question_id = $9
		 RETURNING `+s.answerColumnsToSelect(),
		domain.AnswerStatusEmpty, domain.AnswerStatusWrong, domain.AnswerStatusPending, n(fileID), n(filename), n(textContent), now, userId, questionId,
	), &answer)

	if errors.Is(err, sql.ErrNoRows) {
		return domain.Answer{}, domain.ErrSubmitToNonExistingAnswer
	}
	if err != nil {
		return domain.Answer{}, fmt.Errorf("db: failed to submit answer: %w", err)
	}

	if answer.Status == domain.AnswerStatusCorrect {
		return domain.Answer{}, domain.ErrSubmitToCorrectAnswer
	}
	if answer.Status == domain.AnswerStatusPending && !answer.UpdatedAt.Equal(now) {
		return domain.Answer{}, domain.ErrSubmitToPendingAnswer
	}

	return answer, nil
}

func (s sqlQuestionRepository) GetKnowledgeBars(ctx context.Context, userId int32) ([]domain.KnowledgeBar, error) {
	const query = `
WITH territory_questions AS (
    SELECT 
        iq.territory_id,
        q.question_id,
        q.knowledge_amount
    FROM questions q
    JOIN islands iq ON q.book_id = iq.book_id
)
SELECT 
    tq.territory_id,
    SUM(tq.knowledge_amount) AS total_knowledge,
	SUM(CASE WHEN a.status = $1 THEN tq.knowledge_amount ELSE 0 END) AS achieved_knowledge
FROM territory_questions tq
LEFT JOIN answers a 
    ON tq.question_id = a.question_id
   AND a.user_id = $2
GROUP BY tq.territory_id
ORDER BY tq.territory_id;
`

	rows, err := s.db.QueryContext(ctx, query, domain.AnswerStatusCorrect, userId)
	if err != nil {
		return nil, err
	}
	result := make([]domain.KnowledgeBar, 0, 4)
	for rows.Next() {
		var kb domain.KnowledgeBar
		if err := rows.Scan(&kb.TerritoryID, &kb.Total, &kb.Value); err != nil {
			return nil, err
		}
		result = append(result, kb)
	}

	return result, nil
}

func (s sqlQuestionRepository) HasAnsweredIsland(ctx context.Context, userId int32, islandId string) (bool, error) {
	const query = `
SELECT 
    CASE 
        WHEN COUNT(DISTINCT q.question_id) = COUNT(DISTINCT a.question_id)
        THEN TRUE
        ELSE FALSE
    END AS has_answered_all
FROM islands i
JOIN questions q 
    ON i.book_id = q.book_id
LEFT JOIN answers a
    ON q.question_id = a.question_id
   AND a.user_id = $1
   AND a.status = $2
WHERE i.id = $3 ;
`
	var result bool
	err := s.db.QueryRowContext(ctx, query, userId, domain.AnswerStatusCorrect, islandId).Scan(&result)
	return result, err
}

func (s sqlQuestionRepository) GetQuestion(ctx context.Context, questionId string) (domain.BookQuestion, error) {
	var question domain.BookQuestion
	var rewardSource sql.NullString
	err := s.db.QueryRowContext(ctx, `SELECT question_id, book_id, knowledge_amount, reward_source FROM questions WHERE question_id = $1 ;`,
		questionId).Scan(&question.QuestionID, &question.BookID, &question.KnowledgeAmount, &rewardSource)
	question.RewardSource = rewardSource.String
	if errors.Is(err, sql.ErrNoRows) {
		return question, domain.ErrQuestionNotFound
	}
	return question, err
}

func (s sqlQuestionRepository) CreateCorrection(ctx context.Context, correction domain.Correction) error {
	correction.CreatedAt = time.Now().UTC()
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO corrections (id, question_id, user_id, is_correct, created_at) VALUES ($1, $2, $3, $4, $5)`,
		correction.ID, correction.QuestionId, correction.UserId, correction.IsCorrect, correction.CreatedAt,
	)
	return err
}

func (s sqlQuestionRepository) ApplyCorrection(ctx context.Context, correction domain.Correction, ifBefore time.Time) (bool, error) {
	ifBefore = ifBefore.UTC()
	newStatus := domain.AnswerStatusWrong
	if correction.IsCorrect {
		newStatus = domain.AnswerStatusCorrect
	}
	var postApplyStatus domain.AnswerStatus
	err := s.db.QueryRowContext(ctx,
		`UPDATE answers SET status = CASE WHEN status = $1 THEN $2 ELSE status END WHERE user_id = $3 AND question_id = $4 AND updated_at <= $5 RETURNING status ;`,
		domain.AnswerStatusPending, newStatus, correction.UserId, correction.QuestionId, ifBefore,
	).Scan(&postApplyStatus)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if _, err := s.db.ExecContext(ctx, `UPDATE corrections SET applied = TRUE WHERE id = $1 ;`, correction.ID); err != nil {
		slog.Error("db: failed to set is_applied for correction to true", slog.Any("error", err), slog.String("correction_id", correction.ID))
	}
	if postApplyStatus != newStatus {
		err = domain.ErrAnswerNotPending
	}
	return true, err
}

func (s sqlQuestionRepository) GetUnappliedCorrections(ctx context.Context) (result []domain.Correction, err error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, question_id, user_id, is_correct, created_at FROM corrections WHERE applied = FALSE `,
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		closeErr := rows.Close()
		err = errors.Join(err, closeErr)
	}()
	for rows.Next() {
		var correction domain.Correction
		err := rows.Scan(&correction.ID, &correction.QuestionId, &correction.UserId, &correction.IsCorrect, &correction.CreatedAt)
		if err != nil {
			return nil, err
		}
		result = append(result, correction)
	}
	return result, nil
}
