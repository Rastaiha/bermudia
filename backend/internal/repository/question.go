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
    text TEXT NOT NULL,
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
    feedback TEXT,
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
    status INT4 NOT NULL,
    new_status INT4 NOT NULL,
    feedback TEXT NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_corrections_applied_updated_at ON corrections (status, updated_at);
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
		_, err = tx.ExecContext(ctx, `INSERT INTO questions (question_id, book_id, text, knowledge_amount, reward_source) VALUES ($1, $2, $3, $4, $5) ;`,
			n(q.QuestionID), n(bookId), n(q.Text), q.KnowledgeAmount, n(q.RewardSource),
		)
		if err != nil {
			return fmt.Errorf("insert questions: %w", err)
		}
	}
	return tx.Commit()
}

func (s sqlQuestionRepository) answerColumnsToSelect() string {
	return `user_id, question_id, status, file_id, filename, text_content, feedback, created_at, updated_at`
}

func (s sqlQuestionRepository) scanAnswer(row scannable, answer *domain.Answer) error {
	return row.Scan(&answer.UserID, &answer.QuestionID, &answer.Status, &answer.FileID, &answer.Filename, &answer.TextContent, &answer.Feedback, &answer.CreatedAt, &answer.UpdatedAt)
}

func (s sqlQuestionRepository) GetOrCreateAnswer(ctx context.Context, userId int32, questionID string) (domain.Answer, error) {
	var answer domain.Answer
	now := time.Now().UTC()

	err := s.scanAnswer(s.db.QueryRowContext(ctx,
		`INSERT INTO answers (user_id, question_id, status, created_at, updated_at) 
		 VALUES ($1, $2, $3, $4, $4)
		 ON CONFLICT (user_id, question_id) DO UPDATE SET user_id = EXCLUDED.user_id
		 RETURNING `+s.answerColumnsToSelect(),
		n(userId), n(questionID), domain.AnswerStatusEmpty, now,
	), &answer)

	if err != nil {
		return domain.Answer{}, fmt.Errorf("failed to get or create answer: %w", err)
	}

	return answer, nil
}

func (s sqlQuestionRepository) SubmitAnswer(ctx context.Context, userId int32, questionId, fileID, filename, textContent string) (answer domain.Answer, err error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.Answer{}, fmt.Errorf("start transaction: %w", err)
	}
	defer func() {
		if err != nil {
			err = errors.Join(err, tx.Rollback())
		} else {
			err = tx.Commit()
		}
	}()

	var status domain.AnswerStatus
	err = tx.QueryRowContext(ctx, `SELECT status FROM answers WHERE user_id = $1 AND question_id = $2`,
		userId, questionId).Scan(&status)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Answer{}, domain.ErrSubmitToNonExistingAnswer
	}
	if err != nil {
		return answer, err
	}

	if status == domain.AnswerStatusCorrect || status == domain.AnswerStatusHalfCorrect {
		return domain.Answer{}, domain.ErrSubmitToCorrectAnswer
	}

	if answer.Status == domain.AnswerStatusPending {
		return domain.Answer{}, domain.ErrSubmitToPendingAnswer
	}

	now := time.Now().UTC()

	err = s.scanAnswer(tx.QueryRowContext(ctx,
		`UPDATE answers SET status = $1, file_id = $2, filename = $3, text_content = $4, updated_at = $5
		 WHERE user_id = $6 AND question_id = $7 RETURNING `+s.answerColumnsToSelect(),
		domain.AnswerStatusPending, n(fileID), n(filename), n(textContent), now, userId, questionId,
	), &answer)

	if errors.Is(err, sql.ErrNoRows) {
		return domain.Answer{}, domain.ErrSubmitToNonExistingAnswer
	}
	if err != nil {
		return domain.Answer{}, fmt.Errorf("db: failed to submit answer: %w", err)
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
	SUM(CASE WHEN a.status = $1 THEN tq.knowledge_amount WHEN a.status = $2 THEN tq.knowledge_amount/2 ELSE 0 END) AS achieved_knowledge
FROM territory_questions tq
LEFT JOIN answers a 
    ON tq.question_id = a.question_id
   AND a.user_id = $3
GROUP BY tq.territory_id
ORDER BY tq.territory_id;
`

	rows, err := s.db.QueryContext(ctx, query, domain.AnswerStatusCorrect, domain.AnswerStatusHalfCorrect, userId)
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
   AND (a.status = $2 OR a.status = $3)
WHERE i.id = $4 ;
`
	var result bool
	err := s.db.QueryRowContext(ctx, query, userId, domain.AnswerStatusCorrect, domain.AnswerStatusHalfCorrect, islandId).Scan(&result)
	return result, err
}

func (s sqlQuestionRepository) GetQuestion(ctx context.Context, questionId string) (domain.BookQuestion, error) {
	var question domain.BookQuestion
	var rewardSource sql.NullString
	err := s.db.QueryRowContext(ctx, `SELECT question_id, book_id, text, knowledge_amount, reward_source FROM questions WHERE question_id = $1 ;`,
		questionId).Scan(&question.QuestionID, &question.BookID, &question.Text, &question.KnowledgeAmount, &rewardSource)
	question.RewardSource = rewardSource.String
	if errors.Is(err, sql.ErrNoRows) {
		return question, domain.ErrQuestionNotFound
	}
	return question, err
}

func (s sqlQuestionRepository) CreateCorrection(ctx context.Context, correction domain.Correction) error {
	correction.UpdatedAt = time.Now().UTC()
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO corrections (id, question_id, user_id, status, new_status, feedback, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		correction.ID, correction.QuestionId, correction.UserId, domain.CorrectionStatusDraft, correction.NewStatus, correction.Feedback, correction.UpdatedAt,
	)
	return err
}

func (s sqlQuestionRepository) ApplyCorrection(ctx context.Context, tx domain.Tx, ifBefore time.Time, correction domain.Correction) (domain.Answer, bool, error) {
	if tx == nil {
		tx = s.db
	}
	ifBefore = ifBefore.UTC()
	var answer domain.Answer
	err := s.scanAnswer(tx.QueryRowContext(ctx,
		`UPDATE answers SET
		    status = CASE WHEN status = $1 THEN $2 ELSE status END,
			feedback = CASE WHEN status = $1 THEN $3 ELSE feedback END
            WHERE user_id = $4 AND question_id = $5 AND updated_at <= $6 RETURNING `+s.answerColumnsToSelect(),
		domain.AnswerStatusPending, correction.NewStatus, n(correction.Feedback), correction.UserId, correction.QuestionId, ifBefore,
	), &answer)
	if errors.Is(err, sql.ErrNoRows) {
		return answer, false, nil
	}
	if err != nil {
		return answer, false, err
	}
	if _, err := tx.ExecContext(ctx, `UPDATE corrections SET status = $1 WHERE id = $2 ;`, domain.CorrectionStatusApplied, correction.ID); err != nil {
		slog.Error("db: failed to set correction status to applied", slog.Any("error", err), slog.String("correction_id", correction.ID))
	}
	if answer.Status != correction.NewStatus {
		return answer, false, domain.ErrAnswerNotPending
	}
	return answer, true, nil
}

func (s sqlQuestionRepository) GetUnappliedCorrections(ctx context.Context, before time.Time) (result []domain.Correction, err error) {
	before = before.UTC()
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, question_id, user_id, new_status, feedback, updated_at FROM corrections WHERE status = $1 AND updated_at <= $2 ;`,
		domain.CorrectionStatusPending, before,
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
		err := rows.Scan(&correction.ID, &correction.QuestionId, &correction.UserId, &correction.NewStatus, &correction.Feedback, &correction.UpdatedAt)
		if err != nil {
			return nil, err
		}
		result = append(result, correction)
	}
	return result, nil
}

func (s sqlQuestionRepository) UpdateCorrectionNewStatus(ctx context.Context, id string, newStatus domain.AnswerStatus) error {
	now := time.Now().UTC()
	cmd, err := s.db.ExecContext(ctx, `UPDATE corrections SET new_status = $1, updated_at = $2 WHERE id = $3 AND status = $4 ;`,
		newStatus, now, id, domain.CorrectionStatusDraft)
	if err != nil {
		return err
	}
	rowsAffected, err := cmd.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return domain.ErrAlreadyApplied
	}
	return nil
}

func (s sqlQuestionRepository) UpdateCorrectionFeedback(ctx context.Context, id string, feedback string) (domain.AnswerStatus, error) {
	var newStatus domain.AnswerStatus
	now := time.Now().UTC()
	err := s.db.QueryRowContext(ctx, `UPDATE corrections SET feedback = $1, updated_at = $2 WHERE id = $3 AND status = $4 RETURNING new_status;`,
		feedback, now, id, domain.CorrectionStatusDraft).Scan(&newStatus)
	if errors.Is(err, sql.ErrNoRows) {
		return newStatus, domain.ErrAlreadyApplied
	}
	if err != nil {
		return newStatus, err
	}
	return newStatus, nil
}

func (s sqlQuestionRepository) FinalizeCorrection(ctx context.Context, id string) error {
	now := time.Now().UTC()
	cmd, err := s.db.ExecContext(ctx, `UPDATE corrections SET status = $1, updated_at = $2 WHERE id = $3 AND status = $4 ;`,
		domain.CorrectionStatusPending, now, id, domain.CorrectionStatusDraft)
	if err != nil {
		return err
	}
	rowsAffected, err := cmd.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return domain.ErrAlreadyApplied
	}
	return nil
}
