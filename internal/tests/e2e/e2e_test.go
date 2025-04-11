// internal/integration/e2e_test.go
package integration

import (
	"context"
	"log"
	"net"
	"testing"
	"time"

	"github.com/Chaika-Team/rzd-api/internal/infrastructure/rzd"
	"github.com/Chaika-Team/rzd-api/internal/service"
	transport "github.com/Chaika-Team/rzd-api/internal/transports/grpc"
	pb "github.com/Chaika-Team/rzd-api/internal/transports/grpc/pb"
	"github.com/Chaika-Team/rzd-api/pkg/config"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const testGRPCPort = "50052" // Используйте тестовый порт

// startTestGRPCServer запускает настоящий gRPC-сервер и регистрирует наш обработчик.
func startTestGRPCServer(t *testing.T, svc service.Service) (*grpc.Server, net.Listener) {
	// Создаем эндпоинты и наш gRPC‑обработчик (wrapper)
	eps := transport.MakeEndpoints(svc)
	handler := transport.NewGRPCServer(eps)

	// Создаем настоящий gRPC-сервер
	grpcServer := grpc.NewServer()
	pb.RegisterRzdServiceServer(grpcServer, handler)

	lis, err := net.Listen("tcp", ":"+testGRPCPort)
	require.NoError(t, err)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()
	// Даем время серверу подняться
	time.Sleep(300 * time.Millisecond)
	return grpcServer, lis
}

// newTestGRPCClient создаёт gRPC‑клиента, используя grpc.NewClient вместо DialContext.
func newTestGRPCClient(t *testing.T) pb.RzdServiceClient {
	// Создаем тестовую конфигурацию
	cfg := &config.Config{
		RZD: config.RZD{
			BasePath:    "https://pass.rzd.ru/",
			UserAgent:   "Mozilla/5.0 (compatible; RzdClient/1.0)",
			Language:    "ru",
			Proxy:       "",
			Timeout:     1700,
			RIDLifetime: 300000,
			MaxRetries:  5,
			DebugMode:   false,
		},
		GRPC: config.GRPC{
			Port: testGRPCPort,
		},
	}
	// Создаем реальный клиент RZD (используется в сервисном слое)
	rzdClient, err := rzd.NewRzdClient(&cfg.RZD)
	require.NoError(t, err)

	// Создаем сервисный слой
	svc := service.New(rzdClient)

	// Запускаем тестовый gRPC-сервер
	_, lis := startTestGRPCServer(t, svc)

	// Формируем target с использованием схемы "passthrough"
	target := "passthrough:///" + lis.Addr().String()
	// Используем grpc.NewClient (функция не блокирует подключение)
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)

	t.Cleanup(func() {
		if err := conn.Close(); err != nil {
			t.Logf("failed to close connection: %v", err)
		}
	})
	return pb.NewRzdServiceClient(conn)
}

func TestGetTrainRoutes(t *testing.T) {
	client := newTestGRPCClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := &pb.GetTrainRoutesRequest{
		FromCode:   2004000,
		ToCode:     2000000,
		Direction:  0, // OneWay
		TrainType:  1, // AllTrains
		CheckSeats: false,
		FromDate:   timestamppb.New(time.Now().Add(48 * time.Hour)),
		WithChange: false,
	}

	resp, err := client.GetTrainRoutes(ctx, req)
	require.NoError(t, err)
	require.NotEmpty(t, resp.Routes)

	t.Logf("Получено %d маршрутов", len(resp.Routes))
	for _, route := range resp.Routes {
		t.Logf("Поезд %s, отправление: %s, прибытие: %s", route.TrainNumber, route.Departure.AsTime(), route.Arrival.AsTime())
	}
}

func TestGetTrainCarriages(t *testing.T) {
	client := newTestGRPCClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := &pb.GetTrainCarriagesRequest{
		TrainNumber: "119А",
		Direction:   0,
		FromCode:    2004000,
		FromTime:    timestamppb.New(time.Now().Add(48 * time.Hour)),
		ToCode:      2000000,
	}

	resp, err := client.GetTrainCarriages(ctx, req)
	require.NoError(t, err)
	require.NotEmpty(t, resp.Carriages)

	t.Logf("Получено %d вагонов", len(resp.Carriages))
	for _, car := range resp.Carriages {
		t.Logf("Вагон %s, тип: %s, тариф: %d", car.CarNumber, car.Type, car.Tariff)
	}
}

func TestSearchStation(t *testing.T) {
	client := newTestGRPCClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.SearchStationRequest{
		Query:       "ЧЕБ",
		CompactMode: true,
		Lang:        "ru",
	}

	resp, err := client.SearchStation(ctx, req)
	require.NoError(t, err)
	require.NotEmpty(t, resp.Stations)

	t.Logf("Найдено %d станций", len(resp.Stations))
	for _, s := range resp.Stations {
		t.Logf("Станция: %s, код: %d, уровень: %d, score: %d", s.Name, s.Code, s.Level, s.Score)
	}
}
