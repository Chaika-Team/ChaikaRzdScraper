// pkg/config/config.go
package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config содержит полное конфигурацию приложения
type Config struct {
	RZD  RZD  `env:"RZD"`
	GRPC GRPC `env:"GRPC"`
}

// RZD содержит конфигурацию для клиента RZD
type RZD struct {
	Language    string `env:"RZD_LANGUAGE,default=ru, description=Language of the response"`
	Timeout     int    `env:"RZD_TIMEOUT,default=2000, description=Timeout of retries in milliseconds"`
	MaxRetries  int    `env:"RZD_MAX_RETRIES,default=10, description=Maximum number of retries"`
	RIDLifetime int    `env:"RZD_RID_LIFETIME,default=300000, description=The lifetime of RID in milliseconds"`
	Proxy       string `env:"RZD_PROXY"`
	UserAgent   string `env:"RZD_USER_AGENT,default=Mozilla/5.0 (compatible; RzdClient/1.0)"`
	BasePath    string `env:"RZD_BASE_PATH,default=https://pass.rzd.ru/"`
	DebugMode   bool   `env:"RZD_DEBUG_MODE,default=false"`
}

// GRPC содержит конфигурацию для gRPC сервера
type GRPC struct {
	Port string `env:"GRPC_PORT,default=50051"`
}

// LoadConfig загружает конфигурацию из переменных окружения с использованием cleanenv
func LoadConfig() (*Config, error) {
	cfg := &Config{}
	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("failed to load configuration: %v", err)
	}
	return cfg, nil
}
