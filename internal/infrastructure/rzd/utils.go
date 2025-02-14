package rzd

import "net/http"

// SetHeaders устанавливает заголовки для запросов
func SetHeaders(req *http.Request, client *Client) {
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", client.config.UserAgent)
	req.Header.Set("Referer", client.config.BasePath)
}
