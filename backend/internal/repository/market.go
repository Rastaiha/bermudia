package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Rastaiha/bermudia/internal/domain"
	"time"
)

const marketSchema = `
CREATE TABLE IF NOT EXISTS trade_offers (
	id VARCHAR(255) PRIMARY KEY,
	by INT4 NOT NULL,
	offered TEXT NOT NULL,
	requested TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL,
	deleted_at TIMESTAMP NULL
);

CREATE INDEX IF NOT EXISTS idx_trade_offers_by_created_at ON trade_offers(by, deleted_at, created_at);
CREATE INDEX IF NOT EXISTS idx_trade_offers_created_at ON trade_offers(deleted_at, created_at);
`

type sqlMarketRepository struct {
	db *sql.DB
}

func NewSqlMarketRepository(db *sql.DB) (domain.MarketStore, error) {
	_, err := db.Exec(marketSchema)
	if err != nil {
		return nil, fmt.Errorf("failed to create trade_offers table: %w", err)
	}
	return sqlMarketRepository{db: db}, nil
}

func (s sqlMarketRepository) CreateOffer(ctx context.Context, tx domain.Tx, offer domain.TradeOffer) error {
	if tx == nil {
		tx = s.db
	}

	offeredData, err := json.Marshal(offer.Offered)
	if err != nil {
		return fmt.Errorf("failed to marshal offered cost: %w", err)
	}

	requestedData, err := json.Marshal(offer.Requested)
	if err != nil {
		return fmt.Errorf("failed to marshal requested cost: %w", err)
	}

	_, err = tx.ExecContext(ctx,
		`INSERT INTO trade_offers (id, by, offered, requested, created_at) VALUES ($1, $2, $3, $4, $5)`,
		n(offer.ID), n(offer.By), offeredData, requestedData, offer.CreatedAt,
	)
	return err
}

func (s sqlMarketRepository) DeleteOffer(ctx context.Context, tx domain.Tx, offerId string) error {
	if tx == nil {
		tx = s.db
	}

	cmd, err := tx.ExecContext(ctx,
		`UPDATE trade_offers SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`,
		time.Now().UTC(), offerId,
	)
	if err != nil {
		return err
	}
	affected, err := cmd.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return domain.ErrOfferAlreadyDeleted
	}
	return nil
}

func (s sqlMarketRepository) GetOffer(ctx context.Context, offerId string) (domain.TradeOffer, error) {
	var offer domain.TradeOffer
	var offeredData, requestedData []byte

	err := s.db.QueryRowContext(ctx,
		`SELECT id, by, offered, requested, created_at FROM trade_offers WHERE id = $1 AND deleted_at IS NULL`,
		offerId,
	).Scan(&offer.ID, &offer.By, &offeredData, &requestedData, &offer.CreatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return offer, domain.ErrOfferNotFound
	}

	if err != nil {
		return domain.TradeOffer{}, fmt.Errorf("failed to get trade offer from db: %w", err)
	}

	if err := json.Unmarshal(offeredData, &offer.Offered); err != nil {
		return domain.TradeOffer{}, fmt.Errorf("failed to unmarshal offered cost: %w", err)
	}

	if err := json.Unmarshal(requestedData, &offer.Requested); err != nil {
		return domain.TradeOffer{}, fmt.Errorf("failed to unmarshal requested cost: %w", err)
	}

	return offer, nil
}

func (s sqlMarketRepository) GetOffers(ctx context.Context, filter domain.GetOffersByFilterType, userId int32, offset int, limit int) ([]domain.TradeOffer, error) {
	filterCondition := ""
	switch filter {
	case domain.GetOffersByAll:
	case domain.GetOffersByMe:
		filterCondition = fmt.Sprintf("AND (by = %d)", userId)
	case domain.GetOffersByOthers:
		filterCondition = fmt.Sprintf("AND (by != %d)", userId)
	default:
		return nil, domain.ErrInvalidFilter
	}
	rows, err := s.db.QueryContext(ctx,
		fmt.Sprintf(`SELECT id, by, offered, requested, created_at 
		 FROM trade_offers 
		 WHERE deleted_at IS NULL %s 
		 ORDER BY created_at DESC 
		 LIMIT $1 OFFSET $2`, filterCondition),
		limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query trade offers: %w", err)
	}
	defer rows.Close()

	var offers []domain.TradeOffer
	for rows.Next() {
		var offer domain.TradeOffer
		var offeredData, requestedData []byte

		err := rows.Scan(&offer.ID, &offer.By, &offeredData, &requestedData, &offer.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan trade offer row: %w", err)
		}

		if err := json.Unmarshal(offeredData, &offer.Offered); err != nil {
			return nil, fmt.Errorf("failed to unmarshal offered cost: %w", err)
		}

		if err := json.Unmarshal(requestedData, &offer.Requested); err != nil {
			return nil, fmt.Errorf("failed to unmarshal requested cost: %w", err)
		}

		offers = append(offers, offer)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating trade offer rows: %w", err)
	}

	return offers, nil
}

func (s sqlMarketRepository) GetOffersCountOfUser(ctx context.Context, userId int32) (int, error) {
	var count int
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM trade_offers WHERE by = $1 AND deleted_at IS NULL`, userId).Scan(&count)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("failed to get trade offer count: %w", err)
	}
	return count, nil
}
