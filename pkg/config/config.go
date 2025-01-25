// pkg/config/config.go
package config

import (
	"log"
	"net/url"

	"github.com/ilyakaznacheev/cleanenv"
)

// URLWrapper служит для парсинга URL из строковых переменных окружения
type URLWrapper struct {
	*url.URL
}

// Decode реализует интерфейс cleanenv.Decoder для URLWrapper
func (u *URLWrapper) Decode(value string) error {
	if value == "" {
		u.URL = nil
		return nil
	}
	parsedURL, err := url.Parse(value)
	if err != nil {
		return err
	}
	u.URL = parsedURL
	return nil
}

// Config содержит все конфигурационные параметры сервиса
type Config struct {
	Language  string     `env:"LANGUAGE" env-default:"ru" description:"Language for RZD API"`
	Timeout   float64    `env:"TIMEOUT" env-default:"5.0" description:"HTTP client timeout in seconds"`
	Proxy     URLWrapper `env:"PROXY" env-default:"" description:"Proxy URL"`
	UserAgent string     `env:"USER_AGENT" env-default:"Mozilla/5.0 (compatible; RzdClient/1.0)" description:"User-Agent header for HTTP requests"`
	Referer   string     `env:"REFERER" env-default:"https://pass.rzd.ru/" description:"Referer header for HTTP requests"`
	DebugMode bool       `env:"DEBUG_MODE" env-default:"false" description:"Enable debug mode"`
}

// LoadConfig загружает конфигурацию из переменных окружения с использованием cleanenv
func LoadConfig() *Config {
	cfg := &Config{}
	if err := cleanenv.ReadEnv(cfg); err != nil {
		log.Fatalf("Failed to read environment variables: %v", err)
	}

	return cfg
}
