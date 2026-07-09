package service

import (
	"context"
	"investment-tracker/internal/models"
)

type Repository interface {
	GetByUserID(ctx context.Context, userID int) ([]models.PortfolioItem, error)
}

type Service struct {
	db Repository
}

func New(db Repository) *Service {
	return &Service{db: db}
}

func (s *Service) GetPortfolio(ctx context.Context, userID int) ([]models.PortfolioItem, error) {
	return s.db.GetByUserID(ctx, userID)
}
