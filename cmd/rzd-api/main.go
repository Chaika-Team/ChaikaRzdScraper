// cmd/rzd-api/main.go

package main

import (
	"github.com/Chaika-Team/rzd-api/internal/adapters/config"
	"os"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/gorilla/mux"

	"github.com/Chaika-Team/rzd-api/internal/adapters/http"
	"github.com/Chaika-Team/rzd-api/internal/infrastructure"
	"github.com/Chaika-Team/rzd-api/internal/usecases"
)

func main() {
	// Логгер
	logger := log.NewLogfmtLogger(os.Stdout)
	logger = level.NewFilter(logger, level.AllowDebug())
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	// Загрузка конфигурации
	cfgPath := "./config.yaml"
	cfg := config.GetConfig(logger, cfgPath)

	// Инфраструктура
	httpClient := infrastructure.NewGuzzleHttpClient(cfg)

	// Сервисы
	rzdService := usecases.NewRzdService(httpClient)

	// HTTP Handlers
	router := mux.NewRouter()
	httpHandlers := http.NewHandlers(rzdService, logger)
	httpHandlers.RegisterRoutes(router)

	// Запуск сервера
	addr := cfg.Listen.BindIP + ":" + cfg.Listen.Port
	_ = level.Info(logger).Log("message", "Starting server", "address", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		_ = level.Error(logger).Log("message", "Server failed", "err", err)
		os.Exit(1)
	}
}
