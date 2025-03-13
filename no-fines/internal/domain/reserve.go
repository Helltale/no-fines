package domain

import "context"

type Reserve struct {
	Currency string  `json:"currency"`
	Amount   float64 `json:"amount"`
}

type ReserveRepository interface {
	GetReserves(ctx context.Context) ([]Reserve, error)
	UpdateReserves(ctx context.Context, reserves []Reserve) error
}
