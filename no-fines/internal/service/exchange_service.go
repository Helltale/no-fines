package service

import (
	"context"
	"database/sql"
	"errors"
	"math"
	"time"

	"github.com/Helltale/no-fines/internal/domain"
)

type ExchangeService struct {
	exchangeProviders []domain.ExchangeProvider
	reserveRepo       domain.ReserveRepository
	db                *sql.DB
}

func NewExchangeService(providers []domain.ExchangeProvider, reserveRepo domain.ReserveRepository, db *sql.DB) *ExchangeService {
	return &ExchangeService{
		exchangeProviders: providers,
		reserveRepo:       reserveRepo,
		db:                db,
	}
}

func (s *ExchangeService) GetBestRate(ctx context.Context, pair domain.CurrencyPair) (float64, error) {
	var bestRate float64
	for _, provider := range s.exchangeProviders {
		rates, err := provider.GetRates(pair)
		if err != nil {
			return 0, err
		}
		for _, rate := range rates {
			if rate.Rate > bestRate {
				bestRate = rate.Rate
			}
		}
	}
	if bestRate == 0 {
		return 0, errors.New("no rates available")
	}
	return bestRate, nil
}

func (s *ExchangeService) BuyCurrency(ctx context.Context, currency string, amount float64) error {
	pair := domain.CurrencyPair{
		BaseCurrency:  "RUB",
		QuoteCurrency: currency,
	}
	rate, err := s.GetBestRate(ctx, pair)
	if err != nil {
		return err
	}

	totalCost := amount / rate
	reserves, err := s.reserveRepo.GetReserves(ctx)
	if err != nil {
		return err
	}

	var found bool
	for i, reserve := range reserves {
		if reserve.Currency == "RUB" {
			if reserve.Amount < totalCost {
				return errors.New("insufficient funds to buy currency")
			}
			reserves[i].Amount -= totalCost
			found = true
			break
		}
	}

	if !found {
		return errors.New("base currency not found in reserves")
	}

	for i, reserve := range reserves {
		if reserve.Currency == currency {
			reserves[i].Amount += amount
			break
		}
	}

	return s.reserveRepo.UpdateReserves(ctx, reserves)
}

func (s *ExchangeService) CalculatePNL(startDate, endDate time.Time) (float64, error) {
	query := `
        SELECT SUM(amount * rate * commission) AS pnl
        FROM transactions
        WHERE timestamp BETWEEN $1 AND $2
    `
	var pnl float64
	err := s.db.QueryRow(query, startDate.Unix(), endDate.Unix()).Scan(&pnl)
	if err != nil {
		return 0, err
	}
	return pnl, nil
}

func (s *ExchangeService) FindBestRoute(graph map[string]map[string]float64, start, end string) ([]string, float64, error) {
	distances := make(map[string]float64)
	predecessors := make(map[string]string)
	unvisited := make(map[string]bool)

	for currency := range graph {
		distances[currency] = math.MaxFloat64
		unvisited[currency] = true
	}
	distances[start] = 0

	for len(unvisited) > 0 {
		var current string
		minDist := math.MaxFloat64

		for currency := range unvisited {
			if distances[currency] < minDist {
				minDist = distances[currency]
				current = currency
			}
		}

		delete(unvisited, current)

		if current == end {
			break
		}

		for neighbor, rate := range graph[current] {
			newDist := distances[current] + (1 / rate)
			if newDist < distances[neighbor] {
				distances[neighbor] = newDist
				predecessors[neighbor] = current
			}
		}
	}

	if distances[end] == math.MaxFloat64 {
		return nil, 0, errors.New("no route found")
	}

	path := []string{end}
	for current := end; predecessors[current] != ""; current = predecessors[current] {
		path = append([]string{predecessors[current]}, path...)
	}

	return path, distances[end], nil
}
