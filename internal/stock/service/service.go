package service

import (
	"investment-tracker/internal/models"
	"time"
)

type Repository interface {
	GetStocks() ([]models.Stock, error)
	GetByTicker(t string) (*models.Stock, error)
	CreateStock(stock *models.Stock) error
	GetPriceHistory(ticker string, period time.Time) ([]models.StockPriceLog, error)
}

type Service struct {
	db Repository
}

func New(db Repository) *Service {
	return &Service{db: db}
}

func (s *Service) GetStocks() ([]models.Stock, error) {
	return s.db.GetStocks()
}

func (s *Service) GetByTicker(t string) (*models.Stock, error) {
	return s.db.GetByTicker(t)
}

func (s *Service) CreateStock(stock *models.Stock) error {
	return s.db.CreateStock(stock)
}

func (s *Service) GetPriceHistory(ticker string, period time.Duration) ([]models.StockPriceLog, error) {
	if period == 0 {
		period = 7 * 24 * time.Hour
	}

	startTime := time.Now().Add(-period)

	return s.db.GetPriceHistory(ticker, startTime)
}
