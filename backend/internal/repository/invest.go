package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Rastaiha/bermudia/internal/domain"
)

const (
	investSchema = `
CREATE TABLE IF NOT EXISTS investment_sessions (
    id VARCHAR(255) PRIMARY KEY,
    text TEXT NOT NULL,
    resolved BOOLEAN NOT NULL DEFAULT false,
    end_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create user_investments table
CREATE TABLE IF NOT EXISTS user_investments (
    id SERIAL PRIMARY KEY,
    session_id VARCHAR(255) NOT NULL REFERENCES investment_sessions(id),
    user_id INT4 NOT NULL,
    coin INT4 NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_investment_sessions_active ON investment_sessions (end_at, resolved) ;
CREATE INDEX IF NOT EXISTS idx_user_investments_session_user ON user_investments (session_id, user_id);
`
)

func NewSqlInvestRepository(db *sql.DB) (domain.InvestStore, error) {
	_, err := db.Exec(investSchema)
	if err != nil {
		return nil, fmt.Errorf("failed to create investment sessions table: %w", err)
	}
	return &sqlInvestStore{db: db}, nil
}

type sqlInvestStore struct {
	db *sql.DB
}

func (s *sqlInvestStore) GetSession(ctx context.Context, id string) (*domain.InvestmentSession, error) {
	var session domain.InvestmentSession

	err := s.db.QueryRowContext(ctx, `SELECT id, text, resolved, end_at FROM investment_sessions WHERE id = $1`, id).Scan(
		&session.ID,
		&session.Text,
		&session.Resolved,
		&session.EndAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrInvestSessionNotFound
	}

	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (s *sqlInvestStore) GetActiveSession(ctx context.Context) (*domain.InvestmentSession, error) {
	query := `
		SELECT id, text, resolved, end_at 
		FROM investment_sessions 
		WHERE end_at > CURRENT_TIMESTAMP AND resolved = false
		ORDER BY end_at ASC 
		LIMIT 1
	`

	var session domain.InvestmentSession

	err := s.db.QueryRowContext(ctx, query).Scan(
		&session.ID,
		&session.Text,
		&session.Resolved,
		&session.EndAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // No active session
		}
		return nil, err
	}

	return &session, nil
}

func (s *sqlInvestStore) CreateInvestmentSession(ctx context.Context, tx domain.Tx, session domain.InvestmentSession) error {
	query := `
		INSERT INTO investment_sessions (id, text, resolved, end_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := tx.ExecContext(ctx, query,
		session.ID,
		session.Text,
		session.Resolved,
		session.EndAt,
	)

	return err
}

func (s *sqlInvestStore) AddUserInvestment(ctx context.Context, tx domain.Tx, investment domain.UserInvestment) error {
	query := `
		INSERT INTO user_investments (session_id, user_id, coin)
		VALUES ($1, $2, $3)
	`

	_, err := tx.ExecContext(ctx, query,
		investment.SessionID,
		investment.UserID,
		investment.Coin,
	)

	return err
}

func (s *sqlInvestStore) GetUserInvestments(ctx context.Context, sessionID string, userID int32) ([]domain.UserInvestment, error) {
	query := `
		SELECT session_id, user_id, coin
		FROM user_investments
		WHERE session_id = $1 AND user_id = $2
	`

	rows, err := s.db.QueryContext(ctx, query, sessionID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var investments []domain.UserInvestment
	for rows.Next() {
		var investment domain.UserInvestment
		err := rows.Scan(
			&investment.SessionID,
			&investment.UserID,
			&investment.Coin,
		)
		if err != nil {
			return nil, err
		}
		investments = append(investments, investment)
	}

	return investments, rows.Err()
}

func (s *sqlInvestStore) GetAllUserInvestments(ctx context.Context, sessionID string) ([]domain.UserInvestment, error) {
	query := `
        SELECT session_id, user_id, coin
        FROM user_investments
        WHERE session_id = $1
        ORDER BY user_id ASC
    `

	rows, err := s.db.QueryContext(ctx, query, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var investments []domain.UserInvestment
	for rows.Next() {
		var investment domain.UserInvestment
		err := rows.Scan(
			&investment.SessionID,
			&investment.UserID,
			&investment.Coin,
		)
		if err != nil {
			return nil, err
		}
		investments = append(investments, investment)
	}

	return investments, rows.Err()
}

func (s *sqlInvestStore) MarkResolved(ctx context.Context, tx domain.Tx, sessionID string) error {
	query := `
        UPDATE investment_sessions 
        SET resolved = true, updated_at = CURRENT_TIMESTAMP 
        WHERE id = $1
    `

	result, err := tx.ExecContext(ctx, query, sessionID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrInvestSessionNotFound
	}

	return nil
}
