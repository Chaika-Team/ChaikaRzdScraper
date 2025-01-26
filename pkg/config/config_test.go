// pkg/config/config_test.go
package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	// Устанавливаем переменные окружения
	_ = os.Setenv("RZD_LANGUAGE", "en")
	_ = os.Setenv("RZD_TIMEOUT", "10.5")
	_ = os.Setenv("RZD_PROXY", "http://user:pass@127.0.0.1:8080")
	_ = os.Setenv("RZD_USER_AGENT", "CustomAgent/1.0")
	_ = os.Setenv("RZD_REFERER", "https://example.com/")
	_ = os.Setenv("RZD_DEBUG_MODE", "true")
	_ = os.Setenv("GRPC_PORT", "6000")

	defer func() {
		// Очищаем переменные окружения после теста
		_ = os.Unsetenv("RZD_LANGUAGE")
		_ = os.Unsetenv("RZD_TIMEOUT")
		_ = os.Unsetenv("RZD_PROXY")
		_ = os.Unsetenv("RZD_USER_AGENT")
		_ = os.Unsetenv("RZD_REFERER")
		_ = os.Unsetenv("RZD_DEBUG_MODE")
		_ = os.Unsetenv("GRPC_PORT")
	}()

	cfg := LoadConfig()

	// Проверяем ConfigRZD
	assert.Equal(t, "en", cfg.RZD.Language)
	assert.Equal(t, 10.5, cfg.RZD.Timeout)
	assert.Equal(t, "CustomAgent/1.0", cfg.RZD.UserAgent)
	assert.Equal(t, "https://example.com/", cfg.RZD.BasePath)
	assert.Equal(t, true, cfg.RZD.DebugMode)
	assert.NotNil(t, cfg.RZD.Proxy.URL)
	assert.Equal(t, "http", cfg.RZD.Proxy.Scheme)
	assert.Equal(t, "user:pass@127.0.0.1:8080", cfg.RZD.Proxy.Host)

	// Проверяем ConfigGRPC
	assert.Equal(t, "6000", cfg.GRPC.Port)
}
