// internal/infrastructure/http_client.go

package infrastructure

import (
	"bytes"
	"context"
	"fmt"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/Chaika-Team/rzd-api/internal/config"
)

type HttpClient interface {
	Get(ctx context.Context, url string, params map[string]interface{}) ([]byte, error)
	Post(ctx context.Context, url string, params map[string]interface{}) ([]byte, error)
}

type GuzzleHttpClient struct {
	client *http.Client
	cfg    *config.Config
	logger log.Logger
}

func NewGuzzleHttpClient(cfg *config.Config, logger log.Logger) *GuzzleHttpClient {
	tr := &http.Transport{
		Proxy: nil,
	}

	if cfg.API.Proxy != "" {
		proxyURL, err := url.Parse(cfg.API.Proxy)
		if err == nil {
			tr.Proxy = http.ProxyURL(proxyURL)
		}
	}

	logger = log.With(logger, "component", "http_client")

	return &GuzzleHttpClient{
		client: &http.Client{
			Timeout:   time.Duration(cfg.API.TimeoutSec) * time.Second,
			Transport: tr,
		},
		cfg:    cfg,
		logger: logger,
	}
}

func (c *GuzzleHttpClient) Get(ctx context.Context, urlStr string, params map[string]interface{}) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, err
	}
	// Заполняем query
	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, fmt.Sprintf("%v", v))
	}
	req.URL.RawQuery = q.Encode()

	// Заголовки
	req.Header.Set("User-Agent", c.cfg.API.UserAgent)
	req.Header.Set("Referer", c.cfg.API.Referer)
	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			_ = level.Error(c.logger).Log("msg", "error closing response body", "err", err)
		}
	}(resp.Body)

	return io.ReadAll(resp.Body)
}

func (c *GuzzleHttpClient) Post(ctx context.Context, urlStr string, params map[string]interface{}) ([]byte, error) {
	// form_params в старом PHP => передавал в body form-data?
	// Но на деле Query.php делал 'form_params' => $params, Guzzle отправлял как x-www-form-urlencoded
	// Для идентичности можно сделать JSON, но старый код был url-encoded. Повторим url-encoded:
	formData := url.Values{}
	for k, v := range params {
		formData.Add(k, fmt.Sprintf("%v", v))
	}
	reqBody := bytes.NewBufferString(formData.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, urlStr, reqBody)
	if err != nil {
		return nil, err
	}
	// Заголовки
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", c.cfg.API.UserAgent)
	req.Header.Set("Referer", c.cfg.API.Referer)
	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			_ = level.Error(c.logger).Log("msg", "error closing response body", "err", err)
		}
	}(resp.Body)

	return io.ReadAll(resp.Body)
}
