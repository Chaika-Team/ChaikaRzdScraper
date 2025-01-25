// internal/infrastructure/rzd/rzd_client.go
package rzd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/Chaika-Team/rzd-api/internal/domain"
	"github.com/Chaika-Team/rzd-api/pkg/config"
)

type RzdClient struct {
	BaseURL        string
	SuggestionURL  string
	StationListURL string
	HTTPClient     *http.Client
	Language       string
}

// NewRzdClient инициализирует новый экземпляр RzdClient с конфигурацией
func NewRzdClient(cfg *config.Config) *RzdClient {
	transport := &http.Transport{}

	if cfg.Proxy.URL != nil {
		transport.Proxy = http.ProxyURL(cfg.Proxy.URL)
	}

	if cfg.DebugMode {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	} else {
		log.SetFlags(log.LstdFlags)
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalf("Failed to create CookieJar: %v", err)
	}

	httpClient := &http.Client{
		Timeout:   time.Duration(cfg.Timeout) * time.Second,
		Transport: transport,
		Jar:       jar,
	}

	return &RzdClient{
		BaseURL:        fmt.Sprintf("https://pass.rzd.ru/timetable/public/%s", cfg.Language),
		SuggestionURL:  "https://pass.rzd.ru/suggester",
		StationListURL: "https://pass.rzd.ru/ticket/services/route/basicRoute",
		HTTPClient:     httpClient,
		Language:       cfg.Language,
	}
}

// Метод для выполнения запроса и обработки RID
func (c *RzdClient) executeRequest(req *http.Request) ([]byte, error) {
	log.Printf("Executing request: %s %s", req.Method, req.URL.String())
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		log.Printf("Request failed: %v", err)
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Failed to close response body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		log.Printf("Non-200 response: %d", resp.StatusCode)
		return nil, fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body) //TODO replace deprecated method
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		return nil, err
	}

	var apiResponse map[string]interface{}
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		log.Printf("Failed to unmarshal JSON: %v", err)
		return nil, err
	}

	result, ok := apiResponse["result"].(string)
	if !ok {
		log.Printf("Invalid response structure: no 'result' field")
		return nil, errors.New("invalid response structure: no 'result' field")
	}

	if result == "RID" || result == "REQUEST_ID" {
		rid, ok := apiResponse["rid"].(string)
		if !ok {
			rid, ok = apiResponse["RID"].(string)
			if !ok {
				log.Printf("RID not found in response")
				return nil, errors.New("rid not found in response")
			}
		}

		log.Printf("Received RID: %s", rid)

		// Добавляем rid к параметрам и повторяем запрос
		q := req.URL.Query()
		q.Add("rid", rid)
		req.URL.RawQuery = q.Encode()

		// Повторный запрос с rid
		log.Printf("Retrying request with RID: %s", rid)
		resp, err := c.HTTPClient.Do(req)
		if err != nil {
			log.Printf("Retry request failed: %v", err)
			return nil, err
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Printf("Failed to close retry response body: %v", err)
			}
		}(resp.Body)

		if resp.StatusCode != http.StatusOK {
			log.Printf("Non-200 response on retry: %d", resp.StatusCode)
			return nil, fmt.Errorf("received non-200 response on retry: %d", resp.StatusCode)
		}

		body, err = io.ReadAll(resp.Body) //TODO replace deprecated method
		if err != nil {
			log.Printf("Failed to read retry response body: %v", err)
			return nil, err
		}
	}

	return body, nil
}

// Получение TrainRoutes
func (c *RzdClient) GetTrainRoutes(params domain.GetTrainRoutesParams) ([]domain.TrainRoute, error) {
	layerID := 5827
	reqURL := fmt.Sprintf("%s?layer_id=%d", c.BaseURL, layerID)

	data := url.Values{}
	data.Set("code0", params.Code0)
	data.Set("code1", params.Code1)
	data.Set("dir", fmt.Sprintf("%d", params.Dir))
	data.Set("tfl", fmt.Sprintf("%d", params.Tfl))
	data.Set("checkSeats", fmt.Sprintf("%d", params.CheckSeats))
	data.Set("dt0", params.Dt0)
	data.Set("md", fmt.Sprintf("%d", params.Md))

	req, err := http.NewRequest("POST", reqURL, strings.NewReader(data.Encode()))
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return nil, err
	}

	// Установка заголовков
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; RzdClient/1.0)")
	req.Header.Set("Referer", "https://pass.rzd.ru/")

	responseBody, err := c.executeRequest(req)
	if err != nil {
		log.Printf("Failed to get train routes: %v", err)
		return nil, err
	}

	var apiResponse struct {
		TP []struct {
			List []domain.TrainRoute `json:"list"`
		} `json:"tp"`
	}

	if err := json.Unmarshal(responseBody, &apiResponse); err != nil {
		log.Printf("Failed to unmarshal train routes: %v", err)
		return nil, err
	}

	if len(apiResponse.TP) == 0 {
		log.Printf("No train routes found")
		return nil, errors.New("no train routes found")
	}

	return apiResponse.TP[0].List, nil
}

// Получение TrainRoutesReturn
func (c *RzdClient) GetTrainRoutesReturn(params domain.GetTrainRoutesReturnParams) ([]domain.TrainRoute, []domain.TrainRoute, error) {
	layerID := 5827
	reqURL := fmt.Sprintf("%s?layer_id=%d", c.BaseURL, layerID)

	data := url.Values{}
	data.Set("code0", params.Code0)
	data.Set("code1", params.Code1)
	data.Set("dir", fmt.Sprintf("%d", params.Dir))
	data.Set("tfl", fmt.Sprintf("%d", params.Tfl))
	data.Set("checkSeats", fmt.Sprintf("%d", params.CheckSeats))
	data.Set("dt0", params.Dt0)
	data.Set("dt1", params.Dt1)

	req, err := http.NewRequest("POST", reqURL, strings.NewReader(data.Encode()))
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return nil, nil, err
	}

	// Установка заголовков
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; RzdClient/1.0)")
	req.Header.Set("Referer", "https://pass.rzd.ru/")

	responseBody, err := c.executeRequest(req)
	if err != nil {
		log.Printf("Failed to get train routes return: %v", err)
		return nil, nil, err
	}

	var apiResponse struct {
		TP []struct {
			List []domain.TrainRoute `json:"list"`
		} `json:"tp"`
	}

	if err := json.Unmarshal(responseBody, &apiResponse); err != nil {
		log.Printf("Failed to unmarshal train routes return: %v", err)
		return nil, nil, err
	}

	if len(apiResponse.TP) < 2 {
		log.Printf("Insufficient train routes found for return")
		return nil, nil, errors.New("insufficient train routes found for return")
	}

	return apiResponse.TP[0].List, apiResponse.TP[1].List, nil
}

// Получение TrainCarriages
func (c *RzdClient) GetTrainCarriages(params domain.GetTrainCarriagesParams) (domain.TrainCarriagesResponse, error) {
	layerID := 5764
	reqURL := fmt.Sprintf("%s?layer_id=%d", c.BaseURL, layerID)

	data := url.Values{}
	data.Set("code0", params.Code0)
	data.Set("code1", params.Code1)
	data.Set("tnum0", params.Tnum0)
	data.Set("time0", params.Time0)
	data.Set("dt0", params.Dt0)
	data.Set("dir", fmt.Sprintf("%d", params.Dir))

	req, err := http.NewRequest("POST", reqURL, strings.NewReader(data.Encode()))
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return domain.TrainCarriagesResponse{}, err
	}

	// Установка заголовков
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; RzdClient/1.0)")
	req.Header.Set("Referer", "https://pass.rzd.ru/")

	responseBody, err := c.executeRequest(req)
	if err != nil {
		log.Printf("Failed to get train carriages: %v", err)
		return domain.TrainCarriagesResponse{}, err
	}

	var apiResponse struct {
		Lst []struct {
			Cars           []domain.Carriage `json:"cars"`
			FunctionBlocks []string          `json:"functionBlocks"`
		} `json:"lst"`
		Schemes          []string `json:"schemes"`
		InsuranceCompany []string `json:"insuranceCompany"`
	}

	if err := json.Unmarshal(responseBody, &apiResponse); err != nil {
		log.Printf("Failed to unmarshal train carriages: %v", err)
		return domain.TrainCarriagesResponse{}, err
	}

	if len(apiResponse.Lst) == 0 {
		log.Printf("No carriages found")
		return domain.TrainCarriagesResponse{}, errors.New("no carriages found")
	}

	return domain.TrainCarriagesResponse{
		Cars:           apiResponse.Lst[0].Cars,
		FunctionBlocks: apiResponse.Lst[0].FunctionBlocks,
		Schemes:        apiResponse.Schemes,
		Companies:      apiResponse.InsuranceCompany,
	}, nil
}

// Получение TrainStationList
func (c *RzdClient) GetTrainStationList(params domain.GetTrainStationListParams) (domain.TrainStationListResponse, error) {
	reqURL := c.StationListURL

	data := url.Values{}
	data.Set("trainNumber", params.TrainNumber)
	data.Set("depDate", params.DepDate)
	data.Set("STRUCTURE_ID", "704")

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return domain.TrainStationListResponse{}, err
	}

	req.URL.RawQuery = data.Encode()

	// Установка заголовков
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; RzdClient/1.0)")
	req.Header.Set("Referer", "https://pass.rzd.ru/")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		log.Printf("Request failed: %v", err)
		return domain.TrainStationListResponse{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Failed to close response body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		log.Printf("Non-200 response: %d", resp.StatusCode)
		return domain.TrainStationListResponse{}, fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body) //TODO replace deprecated method
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		return domain.TrainStationListResponse{}, err
	}

	var apiResponse struct {
		Data struct {
			TrainInfo domain.TrainInfo   `json:"trainInfo"`
			Routes    []domain.RouteInfo `json:"routes"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &apiResponse); err != nil {
		log.Printf("Failed to unmarshal train station list: %v", err)
		return domain.TrainStationListResponse{}, err
	}

	return domain.TrainStationListResponse{
		Train:  apiResponse.Data.TrainInfo,
		Routes: apiResponse.Data.Routes,
	}, nil
}

// Получение StationCode
func (c *RzdClient) GetStationCode(params domain.GetStationCodeParams) ([]domain.StationCode, error) {
	reqURL := c.SuggestionURL

	data := url.Values{}
	data.Set("stationNamePart", params.StationNamePart)
	data.Set("compactMode", params.CompactMode)

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return nil, err
	}

	req.URL.RawQuery = data.Encode()

	// Установка заголовков
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; RzdClient/1.0)")
	req.Header.Set("Referer", "https://pass.rzd.ru/")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		log.Printf("Request failed: %v", err)
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Failed to close response body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		log.Printf("Non-200 response: %d", resp.StatusCode)
		return nil, fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body) //TODO replace deprecated method
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		return nil, err
	}

	var apiResponse []struct {
		N string `json:"n"`
		C string `json:"c"`
	}

	if err := json.Unmarshal(body, &apiResponse); err != nil {
		log.Printf("Failed to unmarshal station codes: %v", err)
		return nil, err
	}

	var stations []domain.StationCode
	for _, station := range apiResponse {
		if strings.Contains(strings.ToLower(station.N), strings.ToLower(params.StationNamePart)) {
			stations = append(stations, domain.StationCode{
				Station: station.N,
				Code:    station.C,
			})
		}
	}

	return stations, nil
}
