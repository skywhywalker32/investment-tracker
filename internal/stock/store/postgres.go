package store

import (
	"database/sql"
	"investment-tracker/internal/models"
	"time"
)

const (
	getStocksQuery = `
		SELECT id, ticker, stock_name, stock_price, currency, updated_at, source
		FROM market.stocks
		ORDER BY ticker`
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

func (s *Store) GetStocks() ([]models.Stock, error) {
	rows, err := s.db.Query(getStocksQuery)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var stocks []models.Stock

	for rows.Next() {
		var stock models.Stock

		if err := rows.Scan(
			&stock.ID,
			&stock.Ticker,
			&stock.StockName,
			&stock.StockPrice,
			&stock.Currency,
			&stock.UpdatedAt,
			&stock.Source); err != nil {
			return nil, err
		}

		stocks = append(stocks, stock)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return stocks, nil
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
