// pkg/rzd/rzd_client.go
package rzd

import (
	"bytes"
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

	"github.com/Chaika-Team/rzd-api/internal/utils"

	"github.com/Chaika-Team/rzd-api/internal/domain"
	"github.com/Chaika-Team/rzd-api/internal/infrastructure/rzd/mappers"
	"github.com/Chaika-Team/rzd-api/internal/infrastructure/rzd/schemas"
	"github.com/Chaika-Team/rzd-api/pkg/config"
)

// Client структура клиента
type Client struct {
	config     *config.ConfigRZD
	HTTPClient *http.Client
	Endpoints  Endpoints
	RIDCache   *RIDCache
	mutex      sync.Mutex
}

// NewRzdClient инициализирует новый экземпляр клиента RzdClient с конфигурацией
func NewRzdClient(cfg *config.ConfigRZD) (*Client, error) {
	transport := &http.Transport{}

	if cfg.Proxy != "" {
		proxyURL, err := url.Parse(cfg.Proxy)
		if err != nil {
			return nil, fmt.Errorf("invalid proxy URL: %v", err)
		}
		transport.Proxy = http.ProxyURL(proxyURL)
	}

	// Установка уровня логирования
	if cfg.DebugMode {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	} else {
		log.SetFlags(log.LstdFlags)
	}

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

	// Инициализация клиента
	client := &Client{
		config:     cfg,
		HTTPClient: httpClient,
		Endpoints:  NewEndpoints(cfg.BasePath, cfg.Language),
	}

	return client, nil
}

// executeRequest выполняет HTTP-запрос и обрабатывает ответ, включая обработку RID
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
		err = resp.Body.Close()
		if err != nil {
			return nil, err
		}

		// Разбор JSON-ответа
		var apiResponse map[string]interface{}
		if err := json.Unmarshal(body, &apiResponse); err != nil {
			log.Printf("Failed to unmarshal JSON: %v", err)
			lastError = err
			continue
		}

		result, _ := apiResponse["result"].(string)

		if result == "RID" || result == "REQUEST_ID" {
			rid, err := extractRID(apiResponse)
			if err != nil {
				log.Printf("Failed to extract RID: %v", err)
				lastError = err
				continue
			}
			c.updateRID(rid, time.Duration(c.config.RIDLifetime)*time.Millisecond) // Обновление RID в кэше

			log.Printf("Received RID: %s", rid)

			// Небольшая задержка перед повторным запросом, из c.config.Timeout (int секунд) в time.Duration
			time.Sleep(time.Duration(c.config.Timeout) * time.Millisecond)
			lastError = nil
			continue
		}

		// Проверка на успешный результат
		if result == "OK" {
			// Дополнительная проверка на сообщения об ошибках
			if msg, exists := getErrorMessage(apiResponse); exists {
				log.Printf("API returned error: %s", msg)
				return nil, errors.New(msg)
			}
			return body, nil
		}

		// Обработка других результатов
		log.Printf("Unexpected result field: %s", result)
		lastError = fmt.Errorf("unexpected result field: %s", result)
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
func (c *Client) GetTrainRoutes(params domain.GetTrainRoutesParams) ([]domain.TrainRoute, error) {
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
	setHeaders(req, c)

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

// GetTrainRoutesReturn получает маршруты поездов туда-обратно
func (c *Client) GetTrainRoutesReturn(params domain.GetTrainRoutesReturnParams) ([]domain.TrainRoute, []domain.TrainRoute, error) {
	data := url.Values{}
	data.Set("code0", fmt.Sprintf("%d", params.FromCode))
	data.Set("code1", fmt.Sprintf("%d", params.ToCode))
	data.Set("dir", fmt.Sprintf("%d", params.Direction))
	data.Set("tfl", fmt.Sprintf("%d", params.TrainType))
	data.Set("checkSeats", utils.BoolToString(params.CheckSeats))
	data.Set("dt0", params.FromDate.Format("02.01.2006"))
	data.Set("dt1", params.ToDate.Format("02.01.2006"))

	req, err := http.NewRequest("POST", c.Endpoints.TrainRoutesReturn, strings.NewReader(data.Encode()))
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return nil, nil, err
	}

	// Установка заголовков
	setHeaders(req, c)

	responseBody, err := c.executeRequest(req)
	if err != nil {
		log.Printf("Failed to get train routes return: %v", err)
		return nil, nil, err
	}

	var schemaResp schemas.TrainRouteResponse
	if err := json.Unmarshal(responseBody, &schemaResp); err != nil {
		log.Printf("Failed to unmarshal train routes return: %v", err)
		return nil, nil, err
	}

	if len(schemaResp.TP) < 2 {
		log.Printf("Insufficient train routes found for return")
		return nil, nil, errors.New("insufficient train routes found for return")
	}

	// Маппинг для прямого направления
	forwardRoutes, err := mappers.MapTrainRouteResponse(schemas.TrainRouteResponse{
		Result: schemaResp.Result,
		TP:     []schemas.TP{schemaResp.TP[0]},
	})
	if err != nil {
		return nil, nil, err
	}

	// Маппинг для обратного направления
	backRoutes, err := mappers.MapTrainRouteResponse(schemas.TrainRouteResponse{
		Result: schemaResp.Result,
		TP:     []schemas.TP{schemaResp.TP[1]},
	})
	if err != nil {
		return nil, nil, err
	}

	return forwardRoutes, backRoutes, nil
}

// GetTrainCarriages получает список вагонов выбранного поезда
func (c *Client) GetTrainCarriages(params domain.GetTrainCarriagesParams) (domain.TrainCarriagesResponse, error) {
	data := url.Values{}
	data.Set("code0", fmt.Sprintf("%d", params.FromCode))
	data.Set("code1", fmt.Sprintf("%d", params.FromCode))
	data.Set("tnum0", params.TrainNumber)
	data.Set("time0", params.FromTime.Format("15:04"))
	data.Set("dt0", params.FromDate.Format("02.01.2006"))
	data.Set("dir", fmt.Sprintf("%d", params.Direction))

	req, err := http.NewRequest("POST", c.Endpoints.TrainCarriages, strings.NewReader(data.Encode()))
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return domain.TrainCarriagesResponse{}, err
	}

	// Установка заголовков
	setHeaders(req, c)

	responseBody, err := c.executeRequest(req)
	if err != nil {
		log.Printf("Failed to get train carriages: %v", err)
		return domain.TrainCarriagesResponse{}, err
	}

	var schemaResp schemas.TrainCarriagesResponse
	if err := json.Unmarshal(responseBody, &schemaResp); err != nil {
		log.Printf("Failed to unmarshal train carriages: %v", err)
		return domain.TrainCarriagesResponse{}, err
	}

	// Используем маппер для преобразования схемы в доменные модели
	// domainResp := mappers.MapTrainCarriagesResponse(schemaResp)

	return domain.TrainCarriagesResponse{}, nil
}

// GetTrainStationList получает список станций в маршруте поезда
func (c *Client) GetTrainStationList(params domain.GetTrainStationListParams) (domain.TrainStationListResponse, error) {
	data := url.Values{}
	data.Set("trainNumber", params.TrainNumber)
	data.Set("depDate", params.FromDate.Format("02.01.2006"))
	data.Set("STRUCTURE_ID", fmt.Sprintf("%d", StationsStructureID))

	req, err := http.NewRequest("GET", c.Endpoints.TrainStationList, nil)
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return domain.TrainStationListResponse{}, err
	}

	req.URL.RawQuery = data.Encode()

	// Установка заголовков
	setHeaders(req, c)

	responseBody, err := c.executeRequest(req)
	if err != nil {
		log.Printf("Failed to get train station list: %v", err)
		return domain.TrainStationListResponse{}, err
	}

	var schemaResp schemas.TrainStationListResponse
	if err := json.Unmarshal(responseBody, &schemaResp); err != nil {
		log.Printf("Failed to unmarshal train station list: %v", err)
		return domain.TrainStationListResponse{}, err
	}

	// Используем маппер для преобразования схемы в доменные модели
	//domainResp := mappers.MapTrainStationListResponse(schemaResp)

	return domain.TrainStationListResponse{}, nil
}

// GetStationCode получает список кодов станций по части названия
func (c *Client) GetStationCode(params domain.GetStationCodeParams) ([]domain.StationCode, error) {
	data := url.Values{}
	data.Set("stationNamePart", params.StationNamePart)
	data.Set("compactMode", utils.BoolToString(params.CompactMode))
	req, err := http.NewRequest("GET", c.Endpoints.StationCode, nil)
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return nil, err
	}

	req.URL.RawQuery = data.Encode()

	// Установка заголовков
	setHeaders(req, c)

	responseBody, err := c.executeRequest(req)
	if err != nil {
		log.Printf("Failed to get station codes: %v", err)
		return nil, err
	}

	var schemaResp schemas.StationCodeResponse
	if err := json.Unmarshal(responseBody, &schemaResp); err != nil {
		log.Printf("Failed to unmarshal station codes: %v", err)
		return nil, err
	}

	// Используем маппер для преобразования схемы в доменные модели
	// domainStations := mappers.MapStationCodeResponse(schemaResp)

	return []domain.StationCode{}, nil
}

// setHeaders устанавливает заголовки для запросов
func setHeaders(req *http.Request, client *Client) {
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", client.UserAgent())
	req.Header.Set("Referer", client.config.BasePath)
}

// UserAgent возвращает User-Agent клиента
func (c *Client) UserAgent() string {
	return "Mozilla/5.0 (compatible; RzdClient/1.0)"
}
