// pkg/rzd/rzd_client.go
package rzd

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/Chaika-Team/ChaikaRzdScraper/internal/utils"

	"github.com/Chaika-Team/ChaikaRzdScraper/internal/domain"
	"github.com/Chaika-Team/ChaikaRzdScraper/internal/infrastructure/rzd/mappers"
	"github.com/Chaika-Team/ChaikaRzdScraper/internal/infrastructure/rzd/schemas"
	"github.com/Chaika-Team/ChaikaRzdScraper/pkg/config"
)

// Client структура клиента
type Client struct {
	config     *config.RZD
	HTTPClient *http.Client
	Endpoints  Endpoints
	RIDCache   *RIDCache
	mutex      sync.Mutex
}

// NewRzdClient инициализирует новый экземпляр клиента RzdClient с конфигурацией
func NewRzdClient(cfg *config.RZD) (*Client, error) {
	transport := &http.Transport{}

	if cfg.Proxy != "" {
		proxyURL, err := url.Parse(cfg.Proxy)
		if err != nil {
			return nil, fmt.Errorf("invalid proxy URL: %v", err)
		}
		transport.Proxy = http.ProxyURL(proxyURL)
	}

	// Установка уровня логирования
	// TODO добавить уровень логирования к запросам

	// Создание CookieJar
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create CookieJar: %v", err)
	}

	httpClient := &http.Client{
		Timeout:   time.Duration(cfg.Timeout) * time.Second,
		Transport: transport,
		Jar:       jar,
	}

	endpoints, err := NewEndpoints(cfg.BasePath, cfg.Language)
	if err != nil {
		return nil, fmt.Errorf("failed to create endpoints: %v", err)
	}

	// Инициализация клиента
	client := &Client{
		config:     cfg,
		HTTPClient: httpClient,
		Endpoints:  endpoints,
	}

	return client, nil
}

// executeRequest выполняет HTTP-запрос и обрабатывает ответ, включая обработку RID.
func (c *Client) executeRequest(req *http.Request) ([]byte, error) {
	var lastError error

	for attempt := 1; attempt <= c.config.MaxRetries; attempt++ {
		log.Printf("Executing request: %s %s (Attempt %d)", req.Method, req.URL.String(), attempt)

		// Сохранение тела запроса для повторных попыток
		var reqBodyBytes []byte
		if req.Body != nil {
			reqBodyBytes, _ = io.ReadAll(req.Body)
			req.Body = io.NopCloser(bytes.NewBuffer(reqBodyBytes))
		}

		// Попытка использовать закэшированный RID
		if rid, valid := c.getCachedRID(); valid {
			q := req.URL.Query()
			q.Set("rid", rid)
			req.URL.RawQuery = q.Encode()
			log.Printf("Using cached RID: %s", rid)
		}

		resp, err := c.HTTPClient.Do(req)
		if err != nil {
			log.Printf("Request failed: %v", err)
			lastError = err
			continue
		}

		// Логирование ответа
		respDump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			log.Printf("Failed to dump response: %v", err)
		} else {
			log.Printf("Response dump:\n%s", string(respDump))
		}

		if resp.StatusCode != http.StatusOK {
			log.Printf("Non-200 response: %d", resp.StatusCode)
			err := resp.Body.Close()
			if err != nil {
				return nil, err
			}
			lastError = fmt.Errorf("received non-200 response: %d", resp.StatusCode)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Failed to read response body: %v", err)
			lastError = err
			continue
		}
		if err := resp.Body.Close(); err != nil {
			return nil, err
		}

		// Если ответ начинается с "[", значит это JSON-массив, и проверка поля "result" не требуется.
		trimmedBody := strings.TrimSpace(string(body))
		if strings.HasPrefix(trimmedBody, "[") {
			c.expireRID() // Сброс RID после успешного запроса.
			return body, nil
		}

		// Разбираем JSON-ответ как объект.
		var apiResponse map[string]interface{}
		if err := json.Unmarshal(body, &apiResponse); err != nil {
			log.Printf("Failed to unmarshal JSON: %v", err)
			lastError = err
			continue
		}

		// Если в объекте есть поле "result", работаем с ним.
		result, _ := apiResponse["result"].(string)
		if result == "RID" || result == "REQUEST_ID" {
			rid, err := extractRID(apiResponse)
			if err != nil {
				log.Printf("Failed to extract RID: %v", err)
				lastError = err
				continue
			}
			c.updateRID(rid, time.Duration(c.config.RIDLifetime)*time.Millisecond)
			log.Printf("Received RID: %s", rid)
			// Задержка перед повторным запросом.
			time.Sleep(time.Duration(c.config.Timeout) * time.Millisecond)
			lastError = nil
			continue
		}

		// Если result == "OK", проверяем наличие ошибок.
		if result == "OK" {
			if msg, exists := getErrorMessage(apiResponse); exists {
				log.Printf("API returned error: %s", msg)
				return nil, errors.New(msg)
			}
			c.expireRID()
			return body, nil
		}

		// Обработка других результатов
		log.Printf("Unexpected result field: %s", result)
		lastError = fmt.Errorf("unexpected result field: %s", result)
		time.Sleep(time.Duration(c.config.Timeout) * time.Millisecond)
	}

	return nil, fmt.Errorf("failed after %d attempts: %v", c.config.MaxRetries, lastError)
}

// getErrorMessage извлекает сообщение об ошибке из ответа API, если оно присутствует
func getErrorMessage(apiResponse map[string]interface{}) (string, bool) {
	if tp, ok := apiResponse["tp"].([]interface{}); ok && len(tp) > 0 {
		if tpMap, ok := tp[0].(map[string]interface{}); ok {
			if msgList, ok := tpMap["msgList"].([]interface{}); ok && len(msgList) > 0 {
				if msgMap, ok := msgList[0].(map[string]interface{}); ok {
					if message, ok := msgMap["message"].(string); ok {
						return message, true
					}
				}
			}
		}
	}
	return "", false
}

// GetTrainRoutes получает маршруты поездов в одну точку
func (c *Client) GetTrainRoutes(_ context.Context, params domain.GetTrainRoutesParams) ([]domain.TrainRoute, error) {
	data := url.Values{}
	data.Set("code0", fmt.Sprintf("%d", params.FromCode))
	data.Set("code1", fmt.Sprintf("%d", params.ToCode))
	data.Set("dir", fmt.Sprintf("%d", params.Direction))
	data.Set("tfl", fmt.Sprintf("%d", params.TrainType))
	data.Set("checkSeats", utils.BoolToString(params.CheckSeats))
	data.Set("dt0", params.FromDate.Format("02.01.2006"))
	data.Set("md", utils.BoolToString(params.WithChange))

	req, err := http.NewRequest("POST", c.Endpoints.TrainRoutes, strings.NewReader(data.Encode()))
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return nil, err
	}

	// Установка заголовков
	SetHeaders(req, c)

	responseBody, err := c.executeRequest(req)
	if err != nil {
		log.Printf("Failed to get train routes: %v", err)
		return nil, err
	}

	var schemaResp schemas.TrainRouteResponse
	if err := json.Unmarshal(responseBody, &schemaResp); err != nil {
		log.Printf("Failed to unmarshal train routes: %v", err)
		return nil, err
	}

	// Используем маппер для преобразования схемы в доменные модели
	domainRoutes, err := mappers.MapTrainRouteResponse(schemaResp)
	if err != nil {
		log.Printf("Failed to map train routes: %v", err)
		return nil, err
	}

	return domainRoutes, nil
}

// GetTrainCarriages получает список вагонов выбранного поезда
func (c *Client) GetTrainCarriages(_ context.Context, params domain.GetTrainCarriagesParams) ([]domain.Car, error) {
	data := url.Values{}
	data.Set("code0", fmt.Sprintf("%d", params.FromCode))
	data.Set("code1", fmt.Sprintf("%d", params.ToCode))
	data.Set("tnum0", params.TrainNumber)
	data.Set("time0", params.FromTime.Format("15:04"))
	data.Set("dt0", params.FromTime.Format("02.01.2006"))
	data.Set("dir", fmt.Sprintf("%d", params.Direction))

	req, err := http.NewRequest("POST", c.Endpoints.TrainCarriages, strings.NewReader(data.Encode()))
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return nil, err
	}

	// Установка заголовков
	SetHeaders(req, c)

	responseBody, err := c.executeRequest(req)
	if err != nil {
		log.Printf("Failed to get train carriages: %v", err)
		return nil, err
	}

	var schemaResp schemas.TrainCarriagesResponse
	if err := json.Unmarshal(responseBody, &schemaResp); err != nil {
		log.Printf("Failed to unmarshal train carriages: %v", err)
		return nil, err
	}

	// Используем маппер для преобразования схемы в доменные модели
	domainResp, err := mappers.MapTrainCarriagesResponse(schemaResp)
	if err != nil {
		log.Printf("Failed to map train carriages: %v", err)
		return nil, err
	}

	return domainResp, nil
}

// SearchStation получает список станций, коды которых содержат подстроку запроса.
// Остальные поля ответа игнорируются.
func (c *Client) SearchStation(_ context.Context, params domain.SearchStationParams) ([]domain.Station, error) {
	// Формирование параметров запроса.
	data := url.Values{}
	data.Set("stationNamePart", params.Query)
	data.Set("compactMode", utils.BoolToYesNoLowerCase(params.CompactMode))
	// Добавляем язык из конфигурации
	data.Set("lang", c.config.Language)

	// Создаем GET-запрос к эндпоинту для поиска станций.
	req, err := http.NewRequest("GET", c.Endpoints.StationCode, nil)
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return nil, err
	}
	req.URL.RawQuery = data.Encode()

	// Установка заголовков
	SetHeaders(req, c)

	// Выполняем запрос
	responseBody, err := c.executeRequest(req)
	if err != nil {
		log.Printf("Failed to get station codes: %v", err)
		return nil, err
	}

	// Десериализуем ответ в схему.
	var schemaResp schemas.StationCodeResponse
	if err := json.Unmarshal(responseBody, &schemaResp); err != nil {
		log.Printf("Failed to unmarshal station codes: %v", err)
		return nil, err
	}

	// Используем маппер для преобразования схемы в доменную модель.
	stations, err := mappers.MapStationCodeResponse(schemaResp)
	if err != nil {
		return nil, err
	}

	return stations, nil
}
