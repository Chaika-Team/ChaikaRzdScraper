// cmd/rzd-grpc-service/main.go
package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/Chaika-Team/rzd-api/internal/infrastructure/grpc/pb"
	"github.com/Chaika-Team/rzd-api/internal/infrastructure/rzd"
	"github.com/Chaika-Team/rzd-api/internal/interfaces"
	"github.com/Chaika-Team/rzd-api/internal/usecase"
	"github.com/Chaika-Team/rzd-api/pkg/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	var (
		port string
	)

	flag.StringVar(&port, "port", "50051", "The server port")
	flag.Parse()

	// Загрузка конфигурации
	cfg := config.LoadConfig()

	// Инициализация RzdClient
	rzdClient := rzd.NewRzdClient(cfg)

	// Инициализация UseCase
	routeUseCase := usecase.NewRouteUseCase(rzdClient)
	carriageUseCase := usecase.NewCarriageUseCase(rzdClient)
	stationUseCase := usecase.NewStationUseCase(rzdClient)

	// Создание экземпляра gRPC-сервиса
	rzdServiceServer := interfaces.NewRZDServiceServer(routeUseCase, carriageUseCase, stationUseCase)

	// Создание слушателя
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	// Создание gRPC-сервера
	grpcServer := grpc.NewServer()

	// Регистрация сервиса
	pb.RegisterRZDServiceServer(grpcServer, rzdServiceServer)

	// Включение reflection для удобства тестирования
	reflection.Register(grpcServer)

	log.Printf("RZD gRPC service is running on port %s...", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
