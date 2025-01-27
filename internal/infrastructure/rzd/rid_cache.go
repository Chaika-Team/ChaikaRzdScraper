package rzd

import (
	"errors"
	"fmt"
	"time"
)

// RIDCache кэш RID для запросов к API РЖД
type RIDCache struct {
	RID       string
	ExpiresAt time.Time
}

// extractRID извлекает RID из ответа API
func extractRID(apiResponse map[string]interface{}) (string, error) {
	if rid, ok := apiResponse["RID"]; ok {
		return fmt.Sprintf("%.0f", rid), nil

	}
	if rid, ok := apiResponse["rid"]; ok {
		return fmt.Sprintf("%.0f", rid), nil

	}
	return "", errors.New("rid not found in response")
}

// getCachedRID возвращает RID из кэша
func (c *Client) getCachedRID() (string, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.RIDCache != nil && time.Now().Before(c.RIDCache.ExpiresAt) {
		return c.RIDCache.RID, true
	}
	return "", false
}

// updateRID обновляет RID в кэше c TTL
func (c *Client) updateRID(rid string, ttl time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.RIDCache = &RIDCache{
		RID:       rid,
		ExpiresAt: time.Now().Add(ttl),
	}
}
