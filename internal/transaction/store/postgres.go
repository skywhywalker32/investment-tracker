package store

import (
	"context"
	"database/sql"
	"errors"
	"investment-tracker/internal/models"
)

const (
	insertTransactionQuery = `
		INSERT INTO market.transactions (user_id, ticker_id, operation_type, qty, price, currency)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at`

	upsertPortfolioQuery = `
		INSERT INTO portfolio.portfolios (user_id, ticker_id, qty)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, ticker_id) 
		DO UPDATE SET qty = portfolio.portfolios.qty + EXCLUDED.qty, updated_at = NOW()
		RETURNING qty`

	deleteEmptyPortfolioQuery = `
		DELETE FROM portfolio.portfolios 
		WHERE user_id = $1 AND ticker_id = $2 AND qty = 0`

	getPortfolioQtyQuery = `
		SELECT qty 
		FROM portfolio.portfolios 
		WHERE user_id = $1 AND ticker_id = $2`
)

type Store struct {
	db *sql.DB
}

func New(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetPortfolioQty(ctx context.Context, userID int, tickerID int) (int, error) {
	var qty int
	err := s.db.QueryRowContext(ctx, getPortfolioQtyQuery, userID, tickerID).Scan(&qty)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	}
	return qty, err
}

func (s *Store) CreateTransaction(ctx context.Context, txModel *models.Transaction) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// пишем лог транзакций
	err = tx.QueryRowContext(ctx, insertTransactionQuery,
		txModel.UserID, txModel.TickerID, txModel.OperationType,
		txModel.Qty, txModel.Price, txModel.Currency,
	).Scan(&txModel.ID, &txModel.CreatedAt)
	if err != nil {
		return err
	}

	// считаем изменения для баланса портфеля
	var qtyChange int
	if txModel.OperationType == models.OpBuy {
		qtyChange = txModel.Qty
	} else {
		qtyChange = -txModel.Qty
	}

	// обновляем или создаем запись в портфеле
	var newQty int
	err = tx.QueryRowContext(ctx, upsertPortfolioQuery, txModel.UserID, txModel.TickerID, qtyChange).Scan(&newQty)
	if err != nil {
		return err
	}

	// если продали всё в ноль, то удаляем строку
	if newQty == 0 {
		_, err = tx.ExecContext(ctx, deleteEmptyPortfolioQuery, txModel.UserID, txModel.TickerID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
