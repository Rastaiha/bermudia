package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Rastaiha/bermudia/internal/domain"
	"math/rand"
	"time"
)

const (
	islandsSchema = `
CREATE TABLE IF NOT EXISTS islands (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
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
    pool_id VARCHAR(255) NOT NULL,
    book_id VARCHAR(255) REFERENCES books(id),
    PRIMARY KEY (book_id)
);
CREATE INDEX IF NOT EXISTS idx_books_pool_pool_id ON book_pools (pool_id);
`

	userBooksSchema = `
CREATE TABLE IF NOT EXISTS user_books (
    territory_id VARCHAR(255) NOT NULL,
    island_id VARCHAR(255) NOT NULL,
    user_id INT4 NOT NULL,
    book_id VARCHAR(255) REFERENCES books(id),
    PRIMARY KEY (user_id, island_id)
);
`

	userPortableIslandsSchema = `
CREATE TABLE IF NOT EXISTS user_portable_islands (
    island_id VARCHAR(255) NOT NULL,
    user_id INT4 NOT NULL,
	created_at TIMESTAMP NOT NULL,
    PRIMARY KEY (user_id, island_id)
);
`
)

type sqlIslandRepository struct {
	db *sql.DB
}

func NewSqlIslandRepository(db *sql.DB) (domain.IslandStore, error) {
	_, err := db.Exec(booksSchema)
	if err != nil {
		return nil, fmt.Errorf("create books table: %w", err)
	}
	_, err = db.Exec(islandsSchema)
	if err != nil {
		return nil, fmt.Errorf("create islands table: %w", err)
	}
	_, err = db.Exec(territoryPoolSettingsSchema)
	if err != nil {
		return nil, fmt.Errorf("create territory territory_pool_settings table: %w", err)
	}
	_, err = db.Exec(booksPoolsSchema)
	if err != nil {
		return nil, fmt.Errorf("create book_pools table: %w", err)
	}
	_, err = db.Exec(userBooksSchema)
	if err != nil {
		return nil, fmt.Errorf("create user_books table: %w", err)
	}
	_, err = db.Exec(userPortableIslandsSchema)
	if err != nil {
		return nil, fmt.Errorf("create user_portable_islands table: %w", err)
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
		`INSERT INTO books (id, content) VALUES ($1, $2) ON CONFLICT (id) DO UPDATE SET content = EXCLUDED.content`,
		n(book.ID), c)
	if err != nil {
		return fmt.Errorf("insert book: %w", err)
	}
	return nil
}

func (s sqlIslandRepository) GetBook(ctx context.Context, bookId string) (*domain.Book, error) {
	var content []byte
	err := s.db.QueryRowContext(ctx, `SELECT content FROM books WHERE id = $1`, bookId).Scan(&content)
	if err != nil {
		return nil, fmt.Errorf("get content of book of island: %w", err)
	}
	var result domain.Book
	if err := json.Unmarshal(content, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s sqlIslandRepository) SetIslandHeader(ctx context.Context, header domain.IslandHeader) error {
	if header.BookID != "" && header.FromPool {
		return domain.ErrInvalidIslandHeader
	}
	cmd, err := s.db.ExecContext(ctx,
		`UPDATE islands SET book_id = $1, from_pool = $2 WHERE id = $3 AND territory_id = $4`, n(header.BookID), header.FromPool, header.ID, header.TerritoryID,
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
	rows, err := s.db.QueryContext(ctx, `SELECT id, territory_id, from_pool, book_id FROM islands WHERE territory_id = $1`, territoryId)
	if err != nil {
		return nil, fmt.Errorf("get island headers by territory %q: %w", territoryId, err)
	}
	defer func() {
		err = rows.Close()
	}()
	for rows.Next() {
		var bookId sql.NullString
		var h domain.IslandHeader
		err := rows.Scan(&h.ID, &h.TerritoryID, &h.FromPool, &bookId)
		if err != nil {
			return nil, fmt.Errorf("get island header by territory %q: %w", territoryId, err)
		}
		h.BookID = bookId.String
		result = append(result, h)
	}
	return result, nil
}

func (s sqlIslandRepository) GetIslandHeader(ctx context.Context, islandId string) (domain.IslandHeader, error) {
	var bookId sql.NullString
	var header domain.IslandHeader
	err := s.db.QueryRowContext(ctx, `SELECT id, territory_id, from_pool, book_id FROM islands WHERE id = $1`,
		islandId).Scan(&header.ID, &header.TerritoryID, &header.FromPool, &bookId)
	header.BookID = bookId.String
	if errors.Is(err, sql.ErrNoRows) {
		return header, domain.ErrIslandNotFound
	}
	return header, err
}

func (s sqlIslandRepository) GetIslandHeaderByBookId(ctx context.Context, bookId string) (domain.IslandHeader, error) {
	var bId sql.NullString
	var header domain.IslandHeader
	err := s.db.QueryRowContext(ctx, `SELECT id, territory_id, from_pool, book_id FROM islands WHERE book_id = $1`,
		bookId).Scan(&header.ID, &header.TerritoryID, &header.FromPool, &header.BookID)
	header.BookID = bId.String
	if errors.Is(err, sql.ErrNoRows) {
		return header, domain.ErrIslandNotFound
	}
	return header, err
}

func (s sqlIslandRepository) ReserveIDForTerritory(ctx context.Context, territoryId, islandId, islandName string) error {
	var actualTerritoryId string
	err := s.db.QueryRowContext(ctx,
		`INSERT INTO islands (id, territory_id, name) VALUES ($1, $2, $3) ON CONFLICT (id) DO UPDATE SET name = CASE WHEN $2 = territory_id THEN $3 ELSE name END RETURNING territory_id ;`,
		n(islandId), n(territoryId), n(islandName)).Scan(&actualTerritoryId)
	if err != nil {
		return err
	}
	if actualTerritoryId != territoryId {
		return fmt.Errorf("island_id %q is already taken by territory %q", islandId, actualTerritoryId)
	}
	return nil
}

func (s sqlIslandRepository) GetBookOfIsland(ctx context.Context, islandId string, userId int32) (string, error) {
	var bookId sql.NullString
	var fromPool bool
	err := s.db.QueryRowContext(ctx, `SELECT book_id, from_pool FROM islands WHERE id = $1 `, islandId).
		Scan(&bookId, &fromPool)
	if err != nil {
		return "", fmt.Errorf("get book_id of island: %w", err)
	}
	if bookId.Valid {
		return bookId.String, nil
	}
	if !fromPool {
		return "", domain.ErrEmptyIsland
	}
	err = s.db.QueryRowContext(ctx, `SELECT book_id FROM user_books WHERE island_id = $1 AND user_id = $2`, islandId, userId).Scan(&bookId)
	if errors.Is(err, sql.ErrNoRows) {
		return "", domain.ErrNoBookAssignedFromPool
	}
	if err != nil {
		return "", fmt.Errorf("get book_id of user_books: %w", err)
	}
	return bookId.String, nil
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
	_, err := s.db.ExecContext(ctx, `INSERT INTO book_pools (pool_id, book_id) VALUES ($1, $2) ;`, n(poolId), n(bookId))
	return err
}

func (s sqlIslandRepository) GetPoolOfBook(ctx context.Context, bookId string) (poolId string, found bool, err error) {
	err = s.db.QueryRowContext(ctx, `SELECT pool_id FROM book_pools WHERE book_id = $1`, bookId).Scan(&poolId)
	found = err == nil
	if errors.Is(err, sql.ErrNoRows) {
		err = nil
	}
	return
}

func (s sqlIslandRepository) AssignBookToIslandFromPool(ctx context.Context, territoryId string, islandId string, userId int32) (bookId string, err error) {
	poolCount, err := s.GetTerritoryPoolSettings(ctx, territoryId)
	if err != nil {
		return "", err
	}

	const query = `SELECT bp.pool_id, COUNT(*) FROM user_books ub LEFT JOIN book_pools bp ON ub.book_id = bp.book_id
WHERE user_id = $1 AND territory_id = $2 GROUP BY bp.pool_id`
	rows, err := s.db.QueryContext(ctx, query, userId, territoryId)
	if err != nil {
		return "", err
	}

	for rows.Next() {
		var poolId string
		var count int32
		err := rows.Scan(&poolId, &count)
		if err != nil {
			return "", err
		}
		switch poolId {
		case domain.PoolEasy:
			poolCount.Easy -= count
		case domain.PoolMedium:
			poolCount.Medium -= count
		case domain.PoolHard:
			poolCount.Hard -= count
		}
	}
	err = rows.Close()
	if err != nil {
		return "", err
	}

	chosenPool := ""
	total := poolCount.TotalCount()
	if total <= 0 {
		return "", domain.ErrPoolSettingExhausted
	}
	r := rand.Int31n(total)
	if r < poolCount.Easy {
		chosenPool = domain.PoolEasy
	} else if r >= poolCount.Easy && r < poolCount.Easy+poolCount.Medium {
		chosenPool = domain.PoolMedium
	} else {
		chosenPool = domain.PoolHard
	}

	const query2 = `SELECT bp.book_id FROM book_pools bp WHERE bp.pool_id = $1 AND NOT EXISTS (
    SELECT 1 FROM user_books WHERE user_id = $2 AND book_id = bp.book_id ) ORDER BY RANDOM() LIMIT 1 ;`

	var chosenBookId string
	err = s.db.QueryRowContext(ctx, query2, chosenPool, userId).Scan(&chosenBookId)
	if errors.Is(err, sql.ErrNoRows) {
		return "", domain.ErrBookPoolExhausted
	}
	if err != nil {
		return "", err
	}

	err = s.db.QueryRowContext(ctx, `INSERT INTO user_books (territory_id, island_id, user_id, book_id) VALUES ($1, $2, $3, $4) ON CONFLICT (user_id, island_id) DO UPDATE SET user_id = EXCLUDED.user_id RETURNING book_id ;`,
		n(territoryId), n(islandId), n(userId), n(chosenBookId)).Scan(&bookId)
	return
}

func (s sqlIslandRepository) IsIslandPortable(ctx context.Context, userId int32, islandId string) (bool, error) {
	var result bool
	err := s.db.QueryRowContext(ctx, `SELECT TRUE FROM user_portable_islands WHERE user_id = $1 AND island_id = $2`,
		userId, islandId).Scan(&result)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return result, err
}

func (s sqlIslandRepository) AddPortableIsland(ctx context.Context, userId int32, islandId string) (bool, error) {
	now := time.Now().UTC()
	var actualCreatedAt time.Time
	err := s.db.QueryRowContext(ctx,
		`INSERT INTO user_portable_islands (user_id, island_id, created_at) VALUES ($1, $2, $3)
ON CONFLICT (user_id, island_id) DO UPDATE SET user_id = EXCLUDED.user_id RETURNING created_at;`,
		n(userId), n(islandId), now).Scan(&actualCreatedAt)
	if err != nil {
		return false, err
	}
	return actualCreatedAt.UnixMilli() == now.UnixMilli(), nil
}

func (s sqlIslandRepository) GetPortableIslands(ctx context.Context, userId int32) (result []domain.PortableIsland, err error) {
	result = make([]domain.PortableIsland, 0)
	rows, err := s.db.QueryContext(ctx, `SELECT u.island_id, i.territory_id, i.name FROM user_portable_islands u LEFT JOIN islands i ON u.island_id = i.id WHERE u.user_id = $1 ORDER BY u.created_at;`, userId)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rows.Close()
	}()
	for rows.Next() {
		var p domain.PortableIsland
		err = rows.Scan(&p.IslandID, &p.TerritoryID, &p.Name)
		if err != nil {
			return nil, err
		}
		result = append(result, p)
	}
	return result, rows.Err()
}
