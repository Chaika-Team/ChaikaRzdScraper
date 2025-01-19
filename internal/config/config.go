// internal/config/config.go

package config

import (
	"os"
	"sync"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	API struct {
		BaseURL        string `env:"API_BASE_URL" env-default:"https://pass.rzd.ru/timetable/public/"`
		SuggestionPath string `env:"API_SUGGESTION_PATH" env-default:"https://pass.rzd.ru/suggester"`
		StationListURL string `env:"API_STATION_LIST_URL" env-default:"https://pass.rzd.ru/ticket/services/route/basicRoute"`
		Language       string `env:"API_LANGUAGE" env-default:"ru"`
		UserAgent      string `env:"API_USER_AGENT" env-default:"Mozilla 5"`
		Referer        string `env:"API_REFERER"  env-default:"https://rzd.ru"`
		TimeoutSec     int    `env:"API_TIMEOUT"  env-default:"5"`
		Debug          bool   `env:"API_DEBUG"    env-default:"false"`
		Proxy          string `env:"API_PROXY"    env-default:""`
	} `yaml:"api"`
}

var instance *Config
var once sync.Once

func GetConfig(logger log.Logger, path string) *Config {
	once.Do(func() {
		logger = log.With(logger, "method", "GetConfig")
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
				_ = level.Error(logger).Log("msg", "Failed to read configuration from file", "error", err)
				panic(err)
			}
		}
		_ = level.Info(logger).Log("msg", "Configuration loaded successfully")
	})
	return instance
}
