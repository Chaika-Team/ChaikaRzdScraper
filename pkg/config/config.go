// pkg/config/config.go
package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	RZD  ConfigRZD  `env:"RZD"`
	GRPC ConfigGRPC `env:"GRPC"`
}

type ConfigRZD struct {
	Language  string  `env:"RZD_LANGUAGE,default=ru"`
	Timeout   float64 `env:"RZD_TIMEOUT,default=5"`
	Proxy     string  `env:"RZD_PROXY"`
	UserAgent string  `env:"RZD_USER_AGENT,default=Mozilla/5.0 (compatible; RzdClient/1.0)"`
	BasePath  string  `env:"RZD_BASE_PATH,default=https://pass.rzd.ru/"`
	DebugMode bool    `env:"RZD_DEBUG_MODE,default=false"`
}

type ConfigGRPC struct {
	Port string `env:"GRPC_PORT,default=50051"`
}

// LoadConfig загружает конфигурацию из переменных окружения с использованием cleanenv
func LoadConfig() *Config {
	cfg := &Config{}
	if err := cleanenv.ReadEnv(cfg); err != nil {
		log.Fatalf("Failed to read environment variables: %v", err)
	}

	return cfg
}
