package rzd

import "time"

type RIDCache struct {
	RID       string
	ExpiresAt time.Time
}

func (c *Client) getCachedRID() (string, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.RIDCache != nil && time.Now().Before(c.RIDCache.ExpiresAt) {
		return c.RIDCache.RID, true
	}
	return "", false
}

func (c *Client) updateRID(rid string, ttl time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.RIDCache = &RIDCache{
		RID:       rid,
		ExpiresAt: time.Now().Add(ttl),
	}
}
