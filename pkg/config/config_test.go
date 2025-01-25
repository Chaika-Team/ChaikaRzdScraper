// pkg/config/config_test.go
package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	// Устанавливаем переменные окружения
	_ = os.Setenv("LANGUAGE", "en")
	_ = os.Setenv("TIMEOUT", "10.5")
	_ = os.Setenv("PROXY", "http://user:pass@127.0.0.1:8080")
	_ = os.Setenv("USER_AGENT", "CustomAgent/1.0")
	_ = os.Setenv("REFERER", "https://example.com/")
	_ = os.Setenv("DEBUG_MODE", "true")

	defer func() {
		// Очищаем переменные окружения после теста
		_ = os.Unsetenv("LANGUAGE")
		_ = os.Unsetenv("TIMEOUT")
		_ = os.Unsetenv("PROXY")
		_ = os.Unsetenv("USER_AGENT")
		_ = os.Unsetenv("REFERER")
		_ = os.Unsetenv("DEBUG_MODE")
	}()

	cfg := LoadConfig()

	assert.Equal(t, "en", cfg.Language)
	assert.Equal(t, 10.5, cfg.Timeout)
	assert.Equal(t, "CustomAgent/1.0", cfg.UserAgent)
	assert.Equal(t, "https://example.com/", cfg.Referer)
	assert.Equal(t, true, cfg.DebugMode)
	assert.NotNil(t, cfg.Proxy.URL)
	assert.Equal(t, "http", cfg.Proxy.Scheme)
	assert.Equal(t, "user:pass@127.0.0.1:8080", cfg.Proxy.Host)
}
