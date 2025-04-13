// cmd/rzd-scraper/main.go
package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Chaika-Team/ChaikaRzdScraper/internal/infrastructure/rzd"
	"github.com/Chaika-Team/ChaikaRzdScraper/internal/service"
	"github.com/Chaika-Team/ChaikaRzdScraper/internal/transports/grpc"
	"github.com/Chaika-Team/ChaikaRzdScraper/internal/transports/rabbitmq"
	"github.com/Chaika-Team/ChaikaRzdScraper/pkg/config"
)

// main initializes and runs the RZD integration service. It parses command-line flags for an optional YAML configuration file, sets up context cancellation and signal handling for graceful shutdown, loads the configuration, and initializes the RZD client and service. The function starts the gRPC server and RabbitMQ handler in separate goroutines and waits for a termination signal to gracefully shut down the service.
func main() {
	var (
		configPath string
	)
	// Флаг для опционального пути к файлу конфигурации
	flag.StringVar(&configPath, "config", "", "Путь к YAML файлу конфигурации")
	flag.Parse()

	// Обработка сигналов для graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		cancel()
	}()

	// Загрузка конфигурации (если configPath не пустой, значения берутся из YAML файла)
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	log.Printf("RabbitMQ URL: %s", cfg.RabbitMQ.URL)

	// Инициализация клиента RZD
	client, err := rzd.NewRzdClient(&cfg.RZD)
	if err != nil {
		log.Fatalf("failed to create RZD client: %v", err)
	}

	// Создаем и запускаем сервис
	svc := service.New(client, cfg)
	eps := grpc.MakeEndpoints(svc)
	grpcServer := grpc.NewGRPCServer(eps)

	// Запуск gRPC сервера (ожидается, что функция StartGRPCServer возвращает сервер и listener)
	server, listener, err := grpc.StartGRPCServer(":"+cfg.GRPC.Port, grpcServer)
	if err != nil {
		log.Fatalf("failed to start gRPC server: %v", err)
	}

	// Запуск сервера в отдельной горутине
	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatalf("failed to serve gRPC server: %v", err)
		}
	}()
	log.Printf("gRPC server is running on port %s", cfg.GRPC.Port)

	// Запуск RabbitMQ обработчика в отдельной горутине
	log.Println("Starting RabbitMQ handler...")
	go rabbitmq.StartRabbitMQHandler(svc, cfg.RabbitMQ.URL)

	// Ожидание отмены контекста (сигнала завершения)
	<-ctx.Done()
	log.Println("Shutting down gRPC server...")
	server.GracefulStop()
	log.Println("Server stopped gracefully.")
}
