package main

import (
	"fmt"
	"github.com/go-kit/log"
	"net/http"
	"os"

	"github.com/Chaika-Team/rzd-api/internal/adapters/api"
	"github.com/Chaika-Team/rzd-api/internal/config"
	"github.com/Chaika-Team/rzd-api/internal/infrastructure"
	"github.com/Chaika-Team/rzd-api/internal/usecases"
	"github.com/Chaika-Team/rzd-api/pkg/logger"

	"github.com/gorilla/mux"
)

func main() {
	// Инициализация логгера
	logger := logger.NewLogger("info")
	logger = log.With(logger, "component", "main")

	// Загрузка конфигурации
	cfg := config.GetConfig(logger, "config.yaml")

	// Инициализация клиента API
	httpClient := infrastructure.NewGuzzleHttpClient(cfg, logger)

	// Инициализация RzdAPI
	rzdApi := api.NewRzdAPI(httpClient, cfg, logger)

	// Инициализация сервиса
	rzdService := usecases.NewRzdService(rzdApi)

	// Настройка HTTP Handlers
	handlers := http.NewHandlers(rzdService, logger)

	// Настройка маршрутов
	router := mux.NewRouter()
	handlers.RegisterRoutes(router)

	// Запуск HTTP-сервера
	address := fmt.Sprintf("%s:%s", cfg.HTTP.Host, cfg.HTTP.Port)
	_ = logger.Log("msg", "Starting server", "address", address)

	if err := http.ListenAndServe(address, router); err != nil {
		_ = logger.Log("msg", "Failed to start server", "error", err)
		os.Exit(1)
	}
}
