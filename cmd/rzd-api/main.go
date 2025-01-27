// cmd/rzd-api/main.go
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Chaika-Team/rzd-api/internal/infrastructure/rzd"

	"github.com/Chaika-Team/rzd-api/internal/domain"
	"github.com/Chaika-Team/rzd-api/pkg/config"
)

func main() {
	var (
		port string
	)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		cancel()
	}()

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	// Parse command-line flags
	flag.StringVar(&port, "port", cfg.GRPC.Port, "The gRPC server port")
	flag.Parse()

	cfg.RZD.BasePath = "https://pass.rzd.ru/"
	cfg.RZD.UserAgent = "Mozilla/5.0 (compatible; RzdClient/1.0)"
	cfg.RZD.Language = "ru"
	cfg.RZD.Proxy = ""
	cfg.RZD.Timeout = 1700
	cfg.RZD.RIDLifetime = 300000
	cfg.RZD.MaxRetries = 5
	cfg.RZD.DebugMode = false

	// Тест клиента RZD
	client, err := rzd.NewRzdClient(&cfg.RZD)
	if err != nil {
		log.Fatalf("failed to create RZD client: %v", err)
	}

	params := domain.GetTrainRoutesParams{
		FromCode:   2004000,          // Санкт-Петербург
		ToCode:     2000000,          // Москва
		Direction:  domain.OneWay,    // Только туда
		TrainType:  domain.AllTrains, // Поезда и электрички
		CheckSeats: false,            // Не проверять наличие мест
		FromDate:   time.Now().Add(24 * time.Hour),
		WithChange: false, // Без пересадок
	}

	routes, err := client.GetTrainRoutes(ctx, params)
	if err != nil {
		log.Fatalf("failed to get train routes: %v", err)
	}

	for _, route := range routes {
		fmt.Printf("Поезд %s типа %d из %s в %s отправляется в %s и прибывает в %s\n",
			route.TrainNumber, route.TrainType, route.From.Name, route.To.Name,
			route.Departure.Format("15:04"), route.Arrival.Format("15:04"))
		for _, car := range route.CarTypes {
			fmt.Printf("\tВагон %s %s класса, свободных мест: %d, стоимость: %d руб.\n",
				car.TypeShortLabel, car.Class, car.FreeSeats, car.Tariff)
		}
	}

}
