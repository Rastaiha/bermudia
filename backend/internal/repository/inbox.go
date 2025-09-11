package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/Rastaiha/bermudia/internal/domain"
	"time"
)

const inboxSchema = `
CREATE TABLE IF NOT EXISTS inbox_messages (
    id VARCHAR NOT NULL PRIMARY KEY,
	user_id INT4 NOT NULL,
	content TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_inbox_messages_user_created_at ON inbox_messages(user_id, created_at);
`

type sqlInboxRepository struct {
	db *sql.DB
}

func NewSqlInboxRepository(db *sql.DB) (domain.InboxStore, error) {
	_, err := db.Exec(inboxSchema)
	if err != nil {
		return nil, fmt.Errorf("failed to create inbox_messages table: %w", err)
	}
	return sqlInboxRepository{db: db}, nil
}

func (s sqlInboxRepository) CreateMessage(ctx context.Context, tx domain.Tx, msg domain.InboxMessage) error {
	if tx == nil {
		tx = s.db
	}

	contentData, err := json.Marshal(msg.Content)
	if err != nil {
		return fmt.Errorf("failed to marshal message content: %w", err)
	}

	_, err = tx.ExecContext(ctx,
		`INSERT INTO inbox_messages (id, user_id, content, created_at) VALUES ($1, $2, $3, $4)`,
		n(msg.ID), n(msg.UserID), contentData, msg.CreatedAt,
	)
	return err
}

func (s sqlInboxRepository) GetMessages(ctx context.Context, userId int32, before time.Time, limit int) ([]domain.InboxMessage, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, user_id, content, created_at 
		 FROM inbox_messages 
		 WHERE user_id = $1 AND created_at < $2 
		 ORDER BY created_at DESC 
		 LIMIT $3`,
		userId, before.UTC(), limit,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query inbox messages: %w", err)
	}
	defer rows.Close()

	var messages []domain.InboxMessage
	for rows.Next() {
		var msg domain.InboxMessage
		var contentData []byte

		err := rows.Scan(&msg.ID, &msg.UserID, &contentData, &msg.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan inbox message row: %w", err)
		}

		if err := json.Unmarshal(contentData, &msg.Content); err != nil {
			return nil, fmt.Errorf("failed to unmarshal message content: %w", err)
		}

		messages = append(messages, msg)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating inbox message rows: %w", err)
	}

	return messages, nil
}
