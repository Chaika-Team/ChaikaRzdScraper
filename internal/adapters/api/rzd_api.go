// internal/adapters/api/rzd_api.go

package api

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Chaika-Team/rzd-api/internal/config"
	"github.com/Chaika-Team/rzd-api/internal/domain"
	"github.com/Chaika-Team/rzd-api/internal/infrastructure"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

// RzdAPI реализует интерфейс domain.RzdAPI
type RzdAPI struct {
	client infrastructure.HttpClient
	config *config.Config
	logger log.Logger
}

func NewRzdAPI(client infrastructure.HttpClient, cfg *config.Config, logger log.Logger) domain.RzdAPI {
	logger = log.With(logger, "component", "RzdAPI")
	return &RzdAPI{
		client: client,
		config: cfg,
		logger: logger,
	}
}

// TrainRoutes
// Соответствует старому PHP Api::trainRoutes(array $params): string
func (api *RzdAPI) TrainRoutes(ctx context.Context, p domain.TrainRoutesParams) (string, error) {
	_ = level.Info(api.logger).Log("msg", "Fetching train routes", "params", p)

	// layer_id = 5827
	layer := map[string]interface{}{
		"layer_id": 5827,
	}
	// Сборка финальных параметров (сложение словарей)
	paramsMap := mergeParams(layer, map[string]interface{}{
		"dir":        p.Dir,
		"tfl":        p.Tfl,
		"checkSeats": p.CheckSeats,
		"code0":      p.Code0,
		"code1":      p.Code1,
		"dt0":        p.Dt0.Format("d.m.Y"), // В PHP: format('d.m.Y')
		"md":         p.Md,
	})

	// Выполнение запроса POST (аналог Query->get($this->path, $layer + $params))
	url := fmt.Sprintf("%s%s", api.config.API.BaseURL, api.config.API.Language)
	respBytes, err := api.client.Post(ctx, url, paramsMap)
	if err != nil {
		_ = level.Error(api.logger).Log("msg", "Failed to fetch train routes", "error", err)
		return "", err
	}

	// Парсим ответ в динамический объект, чтобы достать tp[0]->list
	var rawResp struct {
		Tp []struct {
			List interface{} `json:"list"`
		} `json:"tp"`
		Result string `json:"result"`
	}

	if err := json.Unmarshal(respBytes, &rawResp); err != nil {
		_ = level.Error(api.logger).Log("msg", "Failed to parse train routes response", "error", err)
		return "", err
	}
	if rawResp.Result != "" && rawResp.Result != "OK" {
		err := fmt.Errorf("API returned result: %s", rawResp.Result)
		_ = level.Error(api.logger).Log("msg", "API error", "error", err)
		return "", err
	}

	// В старом PHP возвращалось: json_encode($routes->tp[0]->list)
	// rawResp.Tp[0].List - любой тип (массив/список).
	if len(rawResp.Tp) == 0 {
		// Нет tp
		return "[]", nil
	}

	// Финальный JSON
	out, err := json.Marshal(rawResp.Tp[0].List)
	if err != nil {
		return "", err
	}

	_ = level.Info(api.logger).Log("msg", "Fetched train routes successfully")
	return string(out), nil
}

// --------------------- TrainRoutesReturn ---------------------
// Соответствует старому PHP Api::trainRoutesReturn(array $params): string
func (api *RzdAPI) TrainRoutesReturn(ctx context.Context, p domain.TrainRoutesReturnParams) (string, error) {
	_ = level.Info(api.logger).Log("msg", "Fetching round trip routes", "params", p)

	layer := map[string]interface{}{
		"layer_id": 5827,
	}
	paramsMap := mergeParams(layer, map[string]interface{}{
		"dir":        p.Dir,
		"tfl":        p.Tfl,
		"checkSeats": p.CheckSeats,
		"code0":      p.Code0,
		"code1":      p.Code1,
		"dt0":        p.Dt0.Format("d.m.Y"),
		"dt1":        p.Dt1.Format("d.m.Y"),
	})

	url := fmt.Sprintf("%s%s", api.config.API.BaseURL, api.config.API.Language)
	respBytes, err := api.client.Post(ctx, url, paramsMap)
	if err != nil {
		_ = level.Error(api.logger).Log("msg", "Failed to fetch round trip routes", "error", err)
		return "", err
	}

	var rawResp struct {
		Tp []struct {
			List interface{} `json:"list"`
		} `json:"tp"`
		Result string `json:"result"`
	}

	if err := json.Unmarshal(respBytes, &rawResp); err != nil {
		_ = level.Error(api.logger).Log("msg", "Failed to parse round trip routes response", "error", err)
		return "", err
	}
	if rawResp.Result != "" && rawResp.Result != "OK" {
		err := fmt.Errorf("API returned result: %s", rawResp.Result)
		_ = level.Error(api.logger).Log("msg", "API error", "error", err)
		return "", err
	}

	// Старый код возвращал:
	// json_encode([
	//   'forward' => $routes->tp[0]->list,
	//   'back'    => $routes->tp[1]->list,
	// ])
	var forward interface{}
	var back interface{}

	if len(rawResp.Tp) > 0 {
		forward = rawResp.Tp[0].List
	}
	if len(rawResp.Tp) > 1 {
		back = rawResp.Tp[1].List
	}

	outData := map[string]interface{}{
		"forward": forward,
		"back":    back,
	}
	out, err := json.Marshal(outData)
	if err != nil {
		return "", err
	}

	_ = level.Info(api.logger).Log("msg", "Fetched round trip routes successfully")
	return string(out), nil
}

// --------------------- TrainCarriages ---------------------
// Соответствует старому PHP Api::trainCarriages(array $params): string
func (api *RzdAPI) TrainCarriages(ctx context.Context, p domain.TrainCarriagesParams) (string, error) {
	_ = level.Info(api.logger).Log("msg", "Fetching train carriages", "params", p)

	layer := map[string]interface{}{
		"layer_id": 5764,
	}
	paramsMap := mergeParams(layer, map[string]interface{}{
		"dir":   p.Dir,
		"code0": p.Code0,
		"code1": p.Code1,
		"dt0":   p.Dt0.Format("d.m.Y"),
		"time0": p.Time0.Format("15:04"),
		"tnum0": p.Tnum0,
	})

	url := fmt.Sprintf("%s%s", api.config.API.BaseURL, api.config.API.Language)
	respBytes, err := api.client.Post(ctx, url, paramsMap)
	if err != nil {
		_ = level.Error(api.logger).Log("msg", "Failed to fetch train carriages", "error", err)
		return "", err
	}

	var rawResp struct {
		Result string `json:"result"`
		Lst    []struct {
			Cars           interface{} `json:"cars"`
			FunctionBlocks interface{} `json:"functionBlocks"`
		} `json:"lst"`
		Schemes          interface{} `json:"schemes"`
		InsuranceCompany interface{} `json:"insuranceCompany"`
	}

	if err := json.Unmarshal(respBytes, &rawResp); err != nil {
		_ = level.Error(api.logger).Log("msg", "Failed to parse train carriages response", "error", err)
		return "", err
	}
	if rawResp.Result != "" && rawResp.Result != "OK" {
		err := fmt.Errorf("API returned result: %s", rawResp.Result)
		_ = level.Error(api.logger).Log("msg", "API error", "error", err)
		return "", err
	}

	// Старый код:
	// return json_encode([
	//   'cars'           => $carriages->lst[0]->cars ?? null,
	//   'functionBlocks' => $carriages->lst[0]->functionBlocks ?? null,
	//   'schemes'        => $carriages->schemes ?? null,
	//   'companies'      => $carriages->insuranceCompany ?? null,
	// ])
	var cars interface{}
	var functionBlocks interface{}
	if len(rawResp.Lst) > 0 {
		cars = rawResp.Lst[0].Cars
		functionBlocks = rawResp.Lst[0].FunctionBlocks
	}

	outData := map[string]interface{}{
		"cars":           cars,
		"functionBlocks": functionBlocks,
		"schemes":        rawResp.Schemes,
		"companies":      rawResp.InsuranceCompany,
	}
	out, err := json.Marshal(outData)
	if err != nil {
		return "", err
	}

	_ = level.Info(api.logger).Log("msg", "Fetched train carriages successfully")
	return string(out), nil
}

// --------------------- TrainStationList ---------------------
// Соответствует старому PHP Api::trainStationList(array $params): string
func (api *RzdAPI) TrainStationList(ctx context.Context, p domain.TrainStationListParams) (string, error) {
	_ = level.Info(api.logger).Log("msg", "Fetching train station list", "params", p)

	layer := map[string]interface{}{
		"STRUCTURE_ID": 704,
	}
	paramsMap := mergeParams(layer, map[string]interface{}{
		"trainNumber": p.TrainNumber,
		"depDate":     p.DepDate.Format("d.m.Y"),
	})

	// Отличается URL:
	url := api.config.API.StationListURL
	respBytes, err := api.client.Get(ctx, url, paramsMap)
	if err != nil {
		_ = level.Error(api.logger).Log("msg", "Failed to fetch train station list", "error", err)
		return "", err
	}

	// raw объект
	var rawResp struct {
		Data struct {
			TrainInfo interface{} `json:"trainInfo"`
			Routes    interface{} `json:"routes"`
		} `json:"data"`
		Result string `json:"result,omitempty"`
	}

	if err := json.Unmarshal(respBytes, &rawResp); err != nil {
		_ = level.Error(api.logger).Log("msg", "Failed to parse train station list response", "error", err)
		return "", err
	}
	// PHP-версия не проверяла result == OK, поэтому не будем проверять

	// Старый код:
	// return json_encode([
	//   'train'  => $stations->data->trainInfo,
	//   'routes' => $stations->data->routes,
	// ])
	outData := map[string]interface{}{
		"train":  rawResp.Data.TrainInfo,
		"routes": rawResp.Data.Routes,
	}
	out, err := json.Marshal(outData)
	if err != nil {
		return "", err
	}

	_ = level.Info(api.logger).Log("msg", "Fetched train station list successfully")
	return string(out), nil
}

// --------------------- StationCode ---------------------
// Соответствует старому PHP Api::stationCode(array $params): string
func (api *RzdAPI) StationCode(ctx context.Context, p domain.StationCodeParams) (string, error) {
	_ = level.Info(api.logger).Log("msg", "Fetching station codes", "params", p)

	// Старый код: $lang = ['lang' => $this->lang], GET-запрос
	langParam := map[string]interface{}{
		"lang": api.config.API.Language,
	}
	paramsMap := mergeParams(langParam, map[string]interface{}{
		"stationNamePart": p.StationNamePart,
		"compactMode":     p.CompactMode,
	})

	url := api.config.API.SuggestionPath
	respBytes, err := api.client.Get(ctx, url, paramsMap) // GET
	if err != nil {
		_ = level.Error(api.logger).Log("msg", "Failed to fetch station codes", "error", err)
		return "", err
	}

	// JSON => массив объектов { n, c }
	var raw []struct {
		N string `json:"n"`
		C string `json:"c"`
	}

	if err := json.Unmarshal(respBytes, &raw); err != nil {
		_ = level.Error(api.logger).Log("msg", "Failed to parse station codes response", "error", err)
		return "", err
	}

	// Фильтрация как в старом коде: if (mb_stristr($station->n, $params['stationNamePart']))
	// но в Go нет встроенного mb_stristr, имитируем простым Contains (регистр?):
	// Старый код искал подстроку, игнорируя регистр. "stristr" без 'i' - partial case?
	// Упростим, сделаем Contains
	var result []map[string]string
	for _, st := range raw {
		// Если stationNamePart пустой, берём всё
		// или если строка содержится
		if p.StationNamePart == "" || substringCaseInsensitive(st.N, p.StationNamePart) {
			result = append(result, map[string]string{
				"station": st.N,
				"code":    st.C,
			})
		}
	}

	out, err := json.Marshal(result)
	if err != nil {
		return "", err
	}

	_ = level.Info(api.logger).Log("msg", "Fetched station codes successfully", "count", len(result))
	return string(out), nil
}

// --------------------- Вспомогательные функции ---------------------

// mergeParams объединяет два map[string]interface{}
func mergeParams(a, b map[string]interface{}) map[string]interface{} {
	res := make(map[string]interface{}, len(a)+len(b))
	for k, v := range a {
		res[k] = v
	}
	for k, v := range b {
		res[k] = v
	}
	return res
}

// substringCaseInsensitive проверяет, есть ли sub в s без учёта регистра
func substringCaseInsensitive(s, sub string) bool {
	// Упростим через strings.Contains (приводим к lower)
	return containsIgnoreCase(s, sub)
}

func containsIgnoreCase(str, substr string) bool {
	return len(substr) == 0 ||
		len(str) != 0 && // trivial check
			stringContainsFold(str, substr)
}

// stringContainsFold аналог strings.Contains, но без регистра
func stringContainsFold(s, substr string) bool {
	sLower := []runeToLower(s)
	subLower := []runeToLower(substr)
	return naiveContains(sLower, subLower)
}

func runeToLower(s string) []rune {
	var r []rune
	for _, ch := range s {
		// упрощенно
		if ch >= 'A' && ch <= 'Z' {
			ch = ch + ('a' - 'A')
		}
		// можно добавить поддержку юникода при желании
		r = append(r, ch)
	}
	return r
}

func naiveContains(haystack, needle []rune) bool {
	if len(needle) == 0 {
		return true
	}
	if len(haystack) < len(needle) {
		return false
	}
	for i := 0; i+len(needle) <= len(haystack); i++ {
		match := true
		for j := 0; j < len(needle); j++ {
			if haystack[i+j] != needle[j] {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}
