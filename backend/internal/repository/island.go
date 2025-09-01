package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Rastaiha/bermudia/internal/domain"
)

const (
	islandsSchema = `
CREATE TABLE IF NOT EXISTS islands (
    id VARCHAR(255) PRIMARY KEY,
    territory_id VARCHAR(255) NOT NULL,
    book_id VARCHAR(255)
);
CREATE INDEX IF NOT EXISTS idx_islands_territory_id ON islands (territory_id);
CREATE INDEX IF NOT EXISTS idx_islands_book_id ON islands (book_id);
`

	booksSchema = `
CREATE TABLE IF NOT EXISTS books (
    id VARCHAR(255) PRIMARY KEY,
    content TEXT NOT NULL
);
`
)

type sqlIslandRepository struct {
	db *sql.DB
}

func NewSqlIslandRepository(db *sql.DB) (domain.IslandStore, error) {
	_, err := db.Exec(islandsSchema)
	if err != nil {
		return nil, fmt.Errorf("create islands table: %w", err)
	}
	_, err = db.Exec(booksSchema)
	if err != nil {
		return nil, fmt.Errorf("create books table: %w", err)
	}
	return sqlIslandRepository{
		db: db,
	}, nil
}

func (s sqlIslandRepository) SetBook(ctx context.Context, book domain.Book) error {
	c, err := json.Marshal(book)
	if err != nil {
		return err
	}
	_, err = s.db.ExecContext(ctx,
		`INSERT INTO books (id, content) VALUES ($1, $2) ON CONFLICT DO UPDATE SET content = EXCLUDED.content`,
		n(book.ID), c)
	if err != nil {
		return fmt.Errorf("insert book: %w", err)
	}
	return nil
}

func (s sqlIslandRepository) BindBookToIsland(ctx context.Context, islandId string, bookId string) error {
	_, err := s.db.ExecContext(ctx,
		`UPDATE islands SET book_id = $1 WHERE id = $2`,
		n(bookId), n(islandId),
	)
	return err
}

func (s sqlIslandRepository) ReserveIDForTerritory(ctx context.Context, territoryId, islandId string) error {
	var actualTerritoryId string
	err := s.db.QueryRowContext(ctx,
		`INSERT INTO islands (id, territory_id) VALUES ($1, $2) ON CONFLICT DO UPDATE SET id = EXCLUDED.id RETURNING territory_id ;`,
		n(islandId), n(territoryId)).Scan(&actualTerritoryId)
	if err != nil {
		return err
	}
	if actualTerritoryId != territoryId {
		return fmt.Errorf("island_id %q is already taken by territory %q", islandId, actualTerritoryId)
	}
	return nil
}

func (s sqlIslandRepository) GetIslandContent(ctx context.Context, islandId string, userId int32) (*domain.Book, error) {
	var bookId sql.NullString
	err := s.db.QueryRowContext(ctx, `SELECT book_id FROM islands WHERE id = $1 `, islandId).Scan(&bookId)
	if err != nil {
		return nil, fmt.Errorf("get book_id of island: %w", err)
	}
	if !bookId.Valid {
		return &domain.Book{}, nil
	}
	var content []byte
	err = s.db.QueryRowContext(ctx, `SELECT content FROM books WHERE id = $1`, bookId).Scan(&content)
	if err != nil {
		return nil, fmt.Errorf("get content of book of island: %w", err)
	}
	var result domain.Book
	if err := json.Unmarshal(content, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s sqlIslandRepository) GetTerritory(ctx context.Context, id string) (string, error) {
	var territoryId string
	err := s.db.QueryRowContext(ctx, `SELECT territory_id FROM islands WHERE id = $1`, id).Scan(&territoryId)
	if errors.Is(err, sql.ErrNoRows) {
		return "", domain.ErrIslandNotFound
	}
	return territoryId, err
}
