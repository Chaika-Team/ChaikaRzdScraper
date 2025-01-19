// internal/infrastructure/http_client.go

package infrastructure

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/Chaika-Team/rzd-api/internal/config"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type HttpClient interface {
	Get(ctx context.Context, url string, params map[string]interface{}) ([]byte, error)
	Post(ctx context.Context, url string, params map[string]interface{}) ([]byte, error)
}

type GuzzleHttpClient struct {
	client *http.Client
	logger log.Logger
}

func NewGuzzleHttpClient(cfg *config.Config) *GuzzleHttpClient {
	return &GuzzleHttpClient{
		client: &http.Client{
			Timeout: time.Duration(cfg.Storage.MaxConnIdleTime),
		},
	}
}

func (c *GuzzleHttpClient) Get(ctx context.Context, urlStr string, params map[string]interface{}) ([]byte, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	for key, value := range params {
		q.Add(key, value.(string)) // Ensure proper type casting for query parameters
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		_ = level.Error(c.logger).Log("method", "Get", "url", urlStr, "err", err)
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			_ = level.Error(c.logger).Log("msg", "error closing response body", "err", err)
		}
	}(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		_ = level.Error(c.logger).Log("method", "Get", "url", urlStr, "err", err)
		return nil, err
	}

	_ = level.Info(c.logger).Log("method", "Get", "url", urlStr, "status", resp.Status)
	return body, nil
}

func (c *GuzzleHttpClient) Post(ctx context.Context, urlStr string, params map[string]interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", urlStr, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		_ = level.Error(c.logger).Log("method", "Post", "url", urlStr, "err", err)
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			_ = level.Error(c.logger).Log("msg", "error closing response body", "err", err)
		}
	}(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		_ = level.Error(c.logger).Log("method", "Post", "url", urlStr, "err", err)
		return nil, err
	}

	_ = level.Info(c.logger).Log("method", "Post", "url", urlStr, "status", resp.Status)
	return body, nil
}
