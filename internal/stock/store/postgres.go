package store

import (
	"database/sql"
	"investment-tracker/internal/models"
	"time"
)

const (
	getTickersQuery = `
		SELECT ticker 
		FROM market.stocks
		ORDER BY id`

	getByTickerQuery = `
		SELECT id, ticker, stock_name, stock_price, currency, updated_at, source 
		FROM market.stocks 
		WHERE ticker = $1`

	createStockQuery = `
		INSERT INTO market.stocks (ticker, stock_name, stock_price, currency, source) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id, updated_at`

	getPriceHistoryQuery = `
		SELECT id, ticker_id, price, currency, updated_at 
		FROM market.stock_price_logs 
		WHERE ticker_id = (SELECT id FROM market.stocks WHERE ticker = $1) 
		  AND updated_at >= $2
		ORDER BY updated_at`
)

type Store struct {
	db *sql.DB
}

func New(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetTickers() ([]string, error) {
	rows, err := s.db.Query(getTickersQuery)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tickers []string

	for rows.Next() {
		var ticker string

		if err := rows.Scan(&ticker); err != nil {
			return nil, err
		}

		tickers = append(tickers, ticker)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tickers, nil
}

func (s *Store) GetByTicker(t string) (*models.Stock, error) {

	var stock models.Stock

	err := s.db.QueryRow(getByTickerQuery, t).Scan(
		&stock.ID,
		&stock.Ticker,
		&stock.StockName,
		&stock.StockPrice,
		&stock.Currency,
		&stock.UpdatedAt,
		&stock.Source,
	)

	if err != nil {
		return nil, err
	}

	return &stock, nil
}

func (s *Store) CreateStock(stock *models.Stock) error {
	return s.db.QueryRow(
		createStockQuery,
		stock.Ticker, stock.StockName, stock.StockPrice, stock.Currency, stock.Source,
	).Scan(&stock.ID, &stock.UpdatedAt)
}

func (s *Store) GetPriceHistory(ticker string, startTime time.Time) ([]models.StockPriceLog, error) {
	rows, err := s.db.Query(getPriceHistoryQuery, ticker, startTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []models.StockPriceLog
	for rows.Next() {
		var l models.StockPriceLog
		err := rows.Scan(&l.ID, &l.TickerID, &l.Price, &l.Currency, &l.UpdatedAt)
		if err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	return logs, rows.Err()
}
