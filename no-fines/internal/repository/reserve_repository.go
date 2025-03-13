package repository

import (
	"context"
	"database/sql"

	"github.com/Helltale/no-fines/internal/domain"
)

type ReserveRepository struct {
	db *sql.DB
}

func NewReserveRepository(db *sql.DB) *ReserveRepository {
	return &ReserveRepository{db: db}
}

func (r *ReserveRepository) GetReserves(ctx context.Context) ([]domain.Reserve, error) {
	query := `SELECT currency, amount FROM reserves`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reserves []domain.Reserve
	for rows.Next() {
		var reserve domain.Reserve
		if err := rows.Scan(&reserve.Currency, &reserve.Amount); err != nil {
			return nil, err
		}
		reserves = append(reserves, reserve)
	}
	return reserves, nil
}

func (r *ReserveRepository) UpdateReserves(ctx context.Context, reserves []domain.Reserve) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, reserve := range reserves {
		query := `UPDATE reserves SET amount = $1 WHERE currency = $2`
		if _, err := tx.ExecContext(ctx, query, reserve.Amount, reserve.Currency); err != nil {
			return err
		}
	}
	return tx.Commit()
}
