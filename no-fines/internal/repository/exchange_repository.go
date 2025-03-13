package repository

import (
	"context"
	"database/sql"

	"github.com/Helltale/no-fines/internal/domain"
)

type ExchangeRepository struct {
	db *sql.DB
}

func NewExchangeRepository(db *sql.DB) *ExchangeRepository {
	return &ExchangeRepository{db: db}
}

func (r *ExchangeRepository) SaveTransaction(ctx context.Context, tx domain.Transaction) error {
	query := `INSERT INTO transactions (id, user_id, from_currency, to_currency, amount, rate, commission, timestamp, status) 
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.db.ExecContext(ctx, query, tx.ID, tx.UserID, tx.From, tx.To, tx.Amount, tx.Rate, tx.Commission, tx.Timestamp, tx.Status)
	return err
}
