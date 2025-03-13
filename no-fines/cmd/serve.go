package cmd

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/Helltale/no-fines/config"
	pb "github.com/Helltale/no-fines/gen/pb/exchange"
	"github.com/Helltale/no-fines/internal/db"
	"github.com/Helltale/no-fines/internal/domain"
	"github.com/Helltale/no-fines/internal/repository"
	"github.com/Helltale/no-fines/internal/service"
	"github.com/Helltale/no-fines/internal/transport"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "run service",
	Run: func(cmd *cobra.Command, args []string) {

		cfg, err := config.LoadEnv()
		if err != nil {
			log.Fatalln("failed to load env", "error", err)
		}

		postgre, err := db.ConnectWithRetry(cfg, db.ConnectToPostgre)
		if err != nil {
			log.Fatalln("failed to connect to database: ", err)
		}
		defer postgre.Close()

		reserveRepo := repository.NewReserveRepository(postgre)
		mockProvider := &domain.MockProvider{}
		providers := []domain.ExchangeProvider{mockProvider}

		exchangeService := service.NewExchangeService(providers, reserveRepo, postgre)
		reserveService := service.NewReserveService(reserveRepo)

		if err := reserveService.CheckAndReserve(context.Background(), "USD", 100); err != nil {
			log.Fatalf("failed to reserve funds: %v", err)
		}

		// Запускаем gRPC-сервер
		go func() {
			lis, err := net.Listen("tcp", ":"+cfg.APP_GRPC_PORT)
			if err != nil {
				log.Fatalf("failed to listen for gRPC: %v", err)
			}
			server := grpc.NewServer()
			pb.RegisterExchangeServiceServer(server, transport.NewGRPCServer(exchangeService))
			log.Printf("gRPC server listening at %v", lis.Addr())
			if err := server.Serve(lis); err != nil {
				log.Fatalf("failed to serve gRPC: %v", err)
			}
		}()

		// Запускаем HTTP-сервер
		http.Handle("/exchange/rate", transport.NewExchangeHandler(exchangeService))
		log.Printf("HTTP server listening at :%s", cfg.APP_HTTP_PORT)
		if err := http.ListenAndServe(":"+cfg.APP_HTTP_PORT, nil); err != nil {
			log.Fatalf("failed to start HTTP server: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
