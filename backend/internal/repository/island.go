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
    from_pool BOOLEAN NOT NULL DEFAULT FALSE,
    book_id VARCHAR(255) REFERENCES books(id)
);
CREATE INDEX IF NOT EXISTS idx_islands_territory_id ON islands (territory_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_islands_book_id ON islands (book_id);
`

	booksSchema = `
CREATE TABLE IF NOT EXISTS books (
    id VARCHAR(255) PRIMARY KEY,
    content TEXT NOT NULL
);
`

	territoryPoolSettingsSchema = `
CREATE TABLE IF NOT EXISTS territory_pool_settings (
  territory_id VARCHAR(255) PRIMARY KEY,
    easy INT4 NOT NULL,
    medium INT4 NOT NULL,
    hard INT4 NOT NULL
);
`

	booksPoolsSchema = `
CREATE TABLE IF NOT EXISTS book_pools (
    pool_id VARCHAR(255),
    book_id VARCHAR(255) REFERENCES books(id),
    PRIMARY KEY (book_id)
);
CREATE INDEX IF NOT EXISTS idx_books_pool_pool_id ON book_pools (pool_id);
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
	_, err = db.Exec(territoryPoolSettingsSchema)
	if err != nil {
		return nil, fmt.Errorf("create territory territory_pool_settings table: %w", err)
	}
	_, err = db.Exec(booksPoolsSchema)
	if err != nil {
		return nil, fmt.Errorf("create book_pools table: %w", err)
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

func (s sqlIslandRepository) SetIslandHeader(ctx context.Context, territoryId string, header domain.IslandHeader) error {
	if header.BookID != "" && header.FromPool {
		return domain.ErrInvalidIslandHeader
	}
	cmd, err := s.db.ExecContext(ctx,
		`UPDATE islands SET book_id = $1, from_pool = $2 WHERE id = $3 AND territory_id = $4`, n(header.BookID), header.FromPool, header.ID, territoryId,
	)
	if err != nil {
		return fmt.Errorf("update island header: %w", err)
	}
	if i, err := cmd.RowsAffected(); err != nil || i != 1 {
		return domain.ErrInvalidIslandHeader
	}
	return err
}

func (s sqlIslandRepository) GetIslandHeadersByTerritory(ctx context.Context, territoryId string) (result []domain.IslandHeader, err error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, from_pool, book_id FROM islands WHERE territory_id = $1`, territoryId)
	if err != nil {
		return nil, fmt.Errorf("get island headers by territory %q: %w", territoryId, err)
	}
	defer func() {
		err = rows.Close()
	}()
	for rows.Next() {
		var bookId sql.NullString
		var h domain.IslandHeader
		err := rows.Scan(&h.ID, &h.FromPool, &bookId)
		if err != nil {
			return nil, fmt.Errorf("get island header by territory %q: %w", territoryId, err)
		}
		h.BookID = bookId.String
		result = append(result, h)
	}
	return result, nil
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

func (s sqlIslandRepository) SetTerritoryPoolSettings(ctx context.Context, territoryId string, settings domain.TerritoryPoolSettings) error {
	const query = `INSERT INTO territory_pool_settings (territory_id, easy, medium, hard) VALUES ($1, $2, $3, $4)
        ON CONFLICT (territory_id) DO UPDATE SET easy = EXCLUDED.easy, medium = EXCLUDED.medium, hard = EXCLUDED.hard`
	_, err := s.db.ExecContext(ctx, query, territoryId, settings.Easy, settings.Medium, settings.Hard)
	return err
}

func (s sqlIslandRepository) GetTerritoryPoolSettings(ctx context.Context, territoryId string) (domain.TerritoryPoolSettings, error) {
	var settings domain.TerritoryPoolSettings

	err := s.db.QueryRowContext(ctx, `SELECT easy, medium, hard FROM territory_pool_settings WHERE territory_id = $1`, territoryId).
		Scan(&settings.Easy, &settings.Medium, &settings.Hard)

	if err != nil {
		return domain.TerritoryPoolSettings{}, err
	}

	return settings, nil
}

func (s sqlIslandRepository) AddBookToPool(ctx context.Context, poolId string, bookId string) error {
	_, err := s.db.ExecContext(ctx, `INSERT INTO book_pools (pool_id, book_id) VALUES ($1, $2) ;`, poolId, bookId)
	return err
}
