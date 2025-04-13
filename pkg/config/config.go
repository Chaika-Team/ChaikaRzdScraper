// pkg/config/config.go
package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config содержит полное конфигурацию приложения.
type Config struct {
	RZD      RZD      `yaml:"RZD" env:"RZD"`
	GRPC     GRPC     `yaml:"GRPC" env:"GRPC"`
	RabbitMQ RabbitMQ `yaml:"RabbitMQ" env:"RABBITMQ"`
}

// RZD содержит конфигурацию для клиента RZD.
type RZD struct {
	Language    string `yaml:"LANGUAGE" env:"LANGUAGE,default=ru, description=Language of the response"`
	Timeout     int    `yaml:"TIMEOUT" env:"TIMEOUT,default=2000, description=Timeout of retries in milliseconds"`
	MaxRetries  int    `yaml:"MAX_RETRIES" env:"MAX_RETRIES,default=10, description=Maximum number of retries"`
	RIDLifetime int    `yaml:"RID_LIFETIME" env:"RID_LIFETIME,default=300000, description=The lifetime of RID in milliseconds"`
	Proxy       string `yaml:"PROXY" env:"PROXY"`
	UserAgent   string `yaml:"USER_AGENT" env:"USER_AGENT,default=Mozilla/5.0 (compatible; RzdClient/1.0)"`
	BasePath    string `yaml:"BASE_PATH" env:"BASE_PATH,default=https://pass.rzd.ru/"`
	DebugMode   bool   `yaml:"DEBUG_MODE" env:"DEBUG_MODE,default=false"`
}

// GRPC содержит конфигурацию для gRPC сервера.
type GRPC struct {
	Port string `yaml:"PORT" env:"PORT,default=50051"`
}

// RabbitMQ содержит конфигурацию для подключения к RabbitMQ.
type RabbitMQ struct {
	URL string `yaml:"URL" env:"RABBITMQ_URL,default=amqp://guest:guest@localhost:5672/"`
}

// LoadConfig загружает конфигурацию из файла (если передан путь) или из переменных окружения.
// При наличии файла, его значения будут приоритетными.
func LoadConfig(configPath string) (*Config, error) {
	cfg := &Config{}
	if configPath != "" {
		// Если указан файл, считываем конфигурацию из него.
		if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
			return nil, fmt.Errorf("failed to load configuration from file: %v", err)
		}
	} else {
		// Иначе считываем из переменных окружения.
		if err := cleanenv.ReadEnv(cfg); err != nil {
			return nil, fmt.Errorf("failed to load configuration from environment: %v", err)
		}
	}
	return cfg, nil
}
