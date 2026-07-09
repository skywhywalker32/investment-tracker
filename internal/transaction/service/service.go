package service

import (
	"context"
	"errors"
	"investment-tracker/internal/models"
)

var ErrInsufficientStocks = errors.New("insufficient stocks in portfolio to sell")

type Repository interface {
	CreateTransaction(ctx context.Context, txModel *models.Transaction) error
	GetPortfolioQty(ctx context.Context, userID int, tickerID int) (int, error)
}

type Service struct {
	db Repository
}

func New(db Repository) *Service {
	return &Service{db: db}
}

func (s *Service) ExecuteTransaction(ctx context.Context, tx *models.Transaction) (*models.Transaction, error) {
	if tx.Qty <= 0 || tx.Price <= 0 {
		return nil, errors.New("quantity and price must be greater than zero")
	}

	if tx.OperationType == models.OpSell {
		currentQty, err := s.db.GetPortfolioQty(ctx, tx.UserID, tx.TickerID)
		if err != nil {
			return nil, err
		}

		if currentQty < tx.Qty {
			return nil, ErrInsufficientStocks
		}
	}

	err := s.db.CreateTransaction(ctx, tx)
	if err != nil {
		return nil, err
	}

	return tx, nil
}
