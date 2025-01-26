// cmd/rzd-api/main.go
package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/Chaika-Team/rzd-api/internal/infrastructure/rzd"

	"github.com/Chaika-Team/rzd-api/internal/domain"
	"github.com/Chaika-Team/rzd-api/pkg/config"
)

func main() {
	var (
		port string
	)

	// Load configuration
	cfg := config.LoadConfig()

	// Parse command-line flags
	flag.StringVar(&port, "port", cfg.GRPC.Port, "The gRPC server port")
	flag.Parse()

	cfg.RZD.BasePath = "https://pass.rzd.ru/"
	cfg.RZD.UserAgent = "Mozilla/5.0 (compatible; RzdClient/1.0)"
	cfg.RZD.Language = "ru"
	cfg.RZD.Proxy = ""
	cfg.RZD.Timeout = 10
	cfg.RZD.MaxRetries = 10
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

	routes, err := client.GetTrainRoutes(params)
	if err != nil {
		log.Fatalf("failed to get train routes: %v", err)
	}

	for _, route := range routes {
		fmt.Printf("Train %s from %s to %s departs at %s and arrives at %s\n",
			route.TrainNumber, route.From.Name, route.To.Name,
			route.Departure.Format("15:04"), route.Arrival.Format("15:04"))
	}
}
