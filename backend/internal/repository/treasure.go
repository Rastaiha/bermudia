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
	treasuresSchema = `
CREATE TABLE IF NOT EXISTS treasures (
    id VARCHAR(255) PRIMARY KEY,
    book_id VARCHAR(255) NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_treasures_book_id ON treasures (book_id);
`
	userTreasuresSchema = `
CREATE TABLE IF NOT EXISTS user_treasures (
    user_id INT4 NOT NULL,
    treasure_id VARCHAR(255) NOT NULL,
    unlocked BOOLEAN NOT NULL,
    cost TEXT NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    PRIMARY KEY (user_id, treasure_id)
);
CREATE INDEX IF NOT EXISTS idx_user_treasures_user_id ON user_treasures (user_id);
`
)

type sqlTreasureRepository struct {
	db *sql.DB
}

func NewSqlTreasureRepository(db *sql.DB) (domain.TreasureStore, error) {
	_, err := db.Exec(treasuresSchema)
	if err != nil {
		return nil, fmt.Errorf("create treasures table: %w", err)
	}
	_, err = db.Exec(userTreasuresSchema)
	if err != nil {
		return nil, fmt.Errorf("create user_treasures table: %w", err)
	}
	return sqlTreasureRepository{
		db: db,
	}, nil
}

func (s sqlTreasureRepository) BindTreasuresToBook(ctx context.Context, bookId string, treasures []domain.Treasure) (err error) {
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
	_, err = tx.ExecContext(ctx, `DELETE FROM treasures WHERE book_id = $1`, bookId)
	if err != nil {
		return fmt.Errorf("delete treasures: %w", err)
	}
	for _, t := range treasures {
		_, err = tx.ExecContext(ctx, `INSERT INTO treasures (id, book_id) VALUES ($1, $2)`,
			n(t.ID), n(bookId),
		)
		if err != nil {
			return fmt.Errorf("insert treasures: %w", err)
		}
	}
	return tx.Commit()
}

func (s sqlTreasureRepository) userTreasureColumnsToSelect() string {
	return `user_id, treasure_id, unlocked, cost, updated_at`
}

func (s sqlTreasureRepository) scanUserTreasure(row scannable, userTreasure *domain.UserTreasure) error {
	var cost []byte
	err := row.Scan(&userTreasure.UserId, &userTreasure.TreasureID, &userTreasure.Unlocked, &cost, &userTreasure.UpdatedAt)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(cost, &userTreasure.Cost); err != nil {
		return fmt.Errorf("failed to unmarshal cost: %w", err)
	}
	return nil
}

func (s sqlTreasureRepository) GetOrCreateUserTreasure(ctx context.Context, userId int32, treasureId string) (domain.UserTreasure, error) {
	var userTreasure domain.UserTreasure
	now := time.Now().UTC()

	generated := domain.GenerateUserTreasure(userId, treasureId)
	cost, err := json.Marshal(generated.Cost)
	if err != nil {
		return domain.UserTreasure{}, fmt.Errorf("failed to marshal cost: %w", err)
	}

	err = s.scanUserTreasure(s.db.QueryRowContext(ctx,
		`INSERT INTO user_treasures (user_id, treasure_id, unlocked, cost, updated_at) 
		 VALUES ($1, $2, $3, $4, $5)
		 ON CONFLICT (user_id, treasure_id) DO UPDATE SET user_id = EXCLUDED.user_id
		 RETURNING `+s.userTreasureColumnsToSelect(),
		n(userId), n(treasureId), generated.Unlocked, cost, now,
	), &userTreasure)

	if err != nil {
		return domain.UserTreasure{}, fmt.Errorf("failed to get or create user treasure: %w", err)
	}

	return userTreasure, nil
}

func (s sqlTreasureRepository) GetTreasure(ctx context.Context, treasureId string) (domain.Treasure, error) {
	var treasure domain.Treasure
	err := s.db.QueryRowContext(ctx, `SELECT id, book_id FROM treasures WHERE id = $1`,
		treasureId).Scan(&treasure.ID, &treasure.BookID)
	if errors.Is(err, sql.ErrNoRows) {
		return treasure, domain.ErrTreasureNotFound
	}
	return treasure, err
}

func (s sqlTreasureRepository) GetUserTreasure(ctx context.Context, userId int32, treasureId string) (domain.UserTreasure, error) {
	var userTreasure domain.UserTreasure
	err := s.scanUserTreasure(s.db.QueryRowContext(ctx,
		`SELECT `+s.userTreasureColumnsToSelect()+` FROM user_treasures WHERE user_id = $1 AND treasure_id = $2`,
		userId, treasureId,
	), &userTreasure)

	if errors.Is(err, sql.ErrNoRows) {
		return domain.UserTreasure{}, domain.ErrUserTreasureNotFound
	}
	if err != nil {
		return domain.UserTreasure{}, fmt.Errorf("failed to get user treasure from db: %w", err)
	}
	return userTreasure, nil
}

func (s sqlTreasureRepository) UpdateUserTreasure(ctx context.Context, old domain.UserTreasure, updated domain.UserTreasure) error {
	cost, err := json.Marshal(updated.Cost)
	if err != nil {
		return fmt.Errorf("failed to marshal cost: %w", err)
	}
	updated.UpdatedAt = time.Now().UTC()

	cmd, err := s.db.ExecContext(ctx,
		`UPDATE user_treasures
		 SET unlocked = $1, cost = $2, updated_at = $3
		 WHERE user_id = $4 AND treasure_id = $5 AND updated_at = $6`,
		updated.Unlocked, cost, updated.UpdatedAt,
		old.UserId, old.TreasureID, old.UpdatedAt,
	)
	if err != nil {
		return err
	}
	rows, err := cmd.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		// nothing updated -> either user treasure not found, or old didn't match
		return domain.ErrUserTreasureConflict
	}
	return nil
}
