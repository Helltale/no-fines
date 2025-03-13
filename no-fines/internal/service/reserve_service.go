package service

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/Helltale/no-fines/internal/domain"
)

type ReserveService struct {
	reserveRepo domain.ReserveRepository
	mu          sync.Mutex // Для предотвращения race condition
}

func NewReserveService(reserveRepo domain.ReserveRepository) *ReserveService {
	return &ReserveService{
		reserveRepo: reserveRepo,
	}
}

// CheckAndReserve проверяет наличие достаточных резервов и блокирует их.
func (s *ReserveService) CheckAndReserve(ctx context.Context, currency string, amount float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	reserves, err := s.reserveRepo.GetReserves(ctx)
	fmt.Printf("reserve: %v", reserves)
	if err != nil {
		return err
	}

	var found bool
	for i, reserve := range reserves {
		if reserve.Currency == currency {
			if reserve.Amount < amount {
				return errors.New("insufficient reserves")
			}
			reserves[i].Amount -= amount // Блокируем средства
			found = true
			break
		}
	}

	if !found {
		return errors.New("currency not found in reserves")
	}

	return s.reserveRepo.UpdateReserves(ctx, reserves)
}

// ReleaseReserves возвращает заблокированные средства обратно в резервы.
func (s *ReserveService) ReleaseReserves(ctx context.Context, currency string, amount float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	reserves, err := s.reserveRepo.GetReserves(nil)
	if err != nil {
		return err
	}

	for i, reserve := range reserves {
		if reserve.Currency == currency {
			reserves[i].Amount += amount // Возвращаем средства
			return s.reserveRepo.UpdateReserves(ctx, reserves)
		}
	}

	return errors.New("currency not found in reserves")
}

// UpdateReserves обновляет резервы после завершения транзакции.
func (s *ReserveService) UpdateReserves(ctx context.Context, updates []domain.Reserve) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.reserveRepo.UpdateReserves(ctx, updates)
}

func (s *ReserveService) TemporarilyBlockReserves(ctx context.Context, currency string, amount float64, duration time.Duration) error {
	if err := s.CheckAndReserve(ctx, currency, amount); err != nil {
		return err
	}

	go func() {
		<-time.After(duration)
		s.ReleaseReserves(ctx, currency, amount)
	}()
	return nil
}
