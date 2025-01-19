package config

import (
	"os"
	"sync"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/ilyakaznacheev/cleanenv"
)

// Config структура конфигурации приложения.
type Config struct {
	Log struct {
		Level string `env:"LOG_LEVEL" env-default:"info"`
	} `yaml:"log"`
	HTTP struct {
		Host string `env:"HTTP_HOST" env-default:"0.0.0.0"`
		Port string `env:"HTTP_PORT" env-default:"8080"`
	} `yaml:"http"`
	API struct {
		BaseURL    string `env:"API_BASE_URL" env-default:"https://pass.rzd.ru/timetable/public/"`
		Language   string `env:"API_LANGUAGE" env-default:"ru"`
		UserAgent  string `env:"API_USER_AGENT" env-default:"Mozilla/5.0 (compatible; RzdApiWrapper/1.0)"`
		Referer    string `env:"API_REFERER" env-default:"https://rzd.ru"`
		TimeoutSec int    `env:"API_TIMEOUT" env-default:"10"`
	} `yaml:"api"`
}

var instance *Config
var once sync.Once

// GetConfig загружает конфигурацию из переменных окружения или файла.
func GetConfig(logger log.Logger, path string) *Config {
	once.Do(func() {
		logger := log.With(logger, "method", "GetConfig")
		instance = &Config{}

		if _, err := os.Stat(path); os.IsNotExist(err) {
			_ = level.Info(logger).Log("msg", "Config file not found, loading from environment variables")
			if err := cleanenv.ReadEnv(instance); err != nil {
				_ = level.Error(logger).Log("msg", "Failed to read configuration", "error", err)
				panic(err)
			}
		} else {
			_ = level.Info(logger).Log("msg", "Reading configuration from file", "path", path)
			if err := cleanenv.ReadConfig(path, instance); err != nil {
				help, _ := cleanenv.GetDescription(instance, nil)
				_ = level.Error(logger).Log("msg", "Failed to read configuration from file", "error", err, "help", help)
				panic(err)
			}
		}
		_ = level.Info(logger).Log("msg", "Configuration loaded successfully")
	})
	return instance
}
