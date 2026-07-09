package store

import (
	"context"
	"database/sql"
	"investment-tracker/internal/models"
)

const (
	getPortfolioQuery = `
		SELECT s.ticker, s.stock_name, p.qty, s.stock_price AS current_price, s.currency 
		FROM portfolio.portfolios p
		JOIN market.stocks s ON s.id = p.ticker_id
		WHERE p.user_id = $1
		ORDER BY s.ticker`
)

type Store struct {
	db *sql.DB
}

func New(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetByUserID(ctx context.Context, userID int) ([]models.PortfolioItem, error) {
	rows, err := s.db.QueryContext(ctx, getPortfolioQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]models.PortfolioItem, 0)

	for rows.Next() {
		var item models.PortfolioItem
		err := rows.Scan(&item.Ticker, &item.StockName, &item.Qty, &item.CurrentPrice, &item.Currency)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
