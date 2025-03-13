package transport

import (
	"context"
	"log"

	pb "github.com/Helltale/no-fines/gen/pb/exchange"
	"github.com/Helltale/no-fines/internal/domain"
	"github.com/Helltale/no-fines/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCServer struct {
	pb.UnimplementedExchangeServiceServer
	exchangeService *service.ExchangeService
}

func NewGRPCServer(exchangeService *service.ExchangeService) *GRPCServer {
	return &GRPCServer{exchangeService: exchangeService}
}

func (s *GRPCServer) GetExchangeRate(ctx context.Context, req *pb.GetExchangeRateRequest) (*pb.GetExchangeRateResponse, error) {
	pair := domain.CurrencyPair{
		BaseCurrency:  req.BaseCurrency,
		QuoteCurrency: req.QuoteCurrency,
	}
	rate, err := s.exchangeService.GetBestRate(ctx, pair)
	if err != nil {
		log.Println("error getting best rate:", err)
		return nil, status.Error(codes.Internal, "failed to get exchange rate")
	}
	return &pb.GetExchangeRateResponse{Rate: float32(rate)}, nil
}
