// internal/adapters/api/rzd_api.go

package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Chaika-Team/rzd-api/internal/config"
	"time"

	"github.com/Chaika-Team/rzd-api/internal/domain"
	"github.com/Chaika-Team/rzd-api/internal/infrastructure"

	"github.com/google/uuid"
)

type RzdAPI struct {
	client infrastructure.HttpClient
	config config.Config
}

func NewRzdAPI(client infrastructure.HttpClient, config config.Config) domain.RzdAPI {
	return &RzdAPI{
		client: client,
		config: config,
	}
}

func (api *RzdAPI) TrainRoutes(ctx context.Context, params domain.TrainRoutesParams) ([]domain.Trip, error) {
	// Реализация взаимодействия с RZD API для получения маршрутов в одну точку
	// Пример: POST /timetable/public/ru?layer_id=5827&...
	// Аналогично PHP-коду

	// Построение запроса
	requestParams := map[string]interface{}{
		"layer_id":   5827,
		"dir":        params.Dir,
		"tfl":        params.Tfl,
		"checkSeats": params.CheckSeats,
		"code0":      params.Code0,
		"code1":      params.Code1,
		"dt0":        params.Dt0.Format("02.01.2006"),
		"md":         params.Md,
	}

	response, err := api.client.Post(ctx, "https://pass.rzd.ru/timetable/public/"+api.config.Language, requestParams)
	if err != nil {
		return nil, err
	}

	// Парсинг ответа
	var apiResponse struct {
		Result string `json:"result"`
		Tp     []struct {
			List []struct {
				Date0     string `json:"date0"`
				Time0     string `json:"time0"`
				Date1     string `json:"date1"`
				Time1     string `json:"time1"`
				Route0    string `json:"route0"`
				Route1    string `json:"route1"`
				Number    string `json:"number"`
				TimeInWay string `json:"timeInWay"`
				Brand     string `json:"brand"`
				Carrier   string `json:"carrier"`
				Cars      []struct {
					FreeSeats      int     `json:"freeSeats"`
					Itype          string  `json:"itype"`
					ServCls        string  `json:"servCls"`
					Tariff         float64 `json:"tariff"`
					Pt             int     `json:"pt"`
					TypeLoc        string  `json:"typeLoc"`
					Type           string  `json:"type"`
					DisabledPerson bool    `json:"disabledPerson"`
				} `json:"cars"`
			} `json:"list"`
		} `json:"tp"`
	}

	if err := json.Unmarshal(response, &apiResponse); err != nil {
		return nil, err
	}

	if apiResponse.Result != "OK" {
		return nil, fmt.Errorf("API returned result: %s", apiResponse.Result)
	}

	var trips []domain.Trip
	for _, tp := range apiResponse.Tp {
		for _, tripData := range tp.List {
			depTime, err := time.Parse("02.01.2006 15:04", fmt.Sprintf("%s %s", tripData.Date0, tripData.Time0))
			if err != nil {
				return nil, err
			}
			arrTime, err := time.Parse("02.01.2006 15:04", fmt.Sprintf("%s %s", tripData.Date1, tripData.Time1))
			if err != nil {
				return nil, err
			}

			trip := domain.Trip{
				ID:          uuid.New().String(),
				RouteID:     tripData.Route0, // Предполагается, что route0 соответствует RouteID
				TrainNumber: tripData.Number,
				Schedule: domain.Schedule{
					DepartureTime: depTime,
					ArrivalTime:   arrTime,
				},
				DepartureTime: depTime,
				ArrivalTime:   arrTime,
				TrainType:     tripData.Brand,
				Status:        "Запланирован", // Пример статуса
			}

			// Заполнение CarriageConfiguration и других полей при необходимости

			trips = append(trips, trip)
		}
	}

	return trips, nil
}

func (api *RzdAPI) TrainRoutesReturn(ctx context.Context, params domain.TrainRoutesReturnParams) (domain.TripsReturn, error) {
	// Реализация получения маршрутов туда-обратно
	// Аналогично TrainRoutes, но с дополнительными параметрами и разбором ответа

	requestParams := map[string]interface{}{
		"layer_id":   5827,
		"dir":        params.Dir,
		"tfl":        params.Tfl,
		"checkSeats": params.CheckSeats,
		"code0":      params.Code0,
		"code1":      params.Code1,
		"dt0":        params.Dt0.Format("02.01.2006"),
		"dt1":        params.Dt1.Format("02.01.2006"),
	}

	response, err := api.client.Post(ctx, "https://pass.rzd.ru/timetable/public/"+api.config.Language, requestParams)
	if err != nil {
		return domain.TripsReturn{}, err
	}

	var apiResponse struct {
		Result string `json:"result"`
		Tp     []struct {
			List []struct {
				// Поля аналогичны TrainRoutes
				Date0     string `json:"date0"`
				Time0     string `json:"time0"`
				Date1     string `json:"date1"`
				Time1     string `json:"time1"`
				Route0    string `json:"route0"`
				Route1    string `json:"route1"`
				Number    string `json:"number"`
				TimeInWay string `json:"timeInWay"`
				Brand     string `json:"brand"`
				Carrier   string `json:"carrier"`
				Cars      []struct {
					FreeSeats      int     `json:"freeSeats"`
					Itype          string  `json:"itype"`
					ServCls        string  `json:"servCls"`
					Tariff         float64 `json:"tariff"`
					Pt             int     `json:"pt"`
					TypeLoc        string  `json:"typeLoc"`
					Type           string  `json:"type"`
					DisabledPerson bool    `json:"disabledPerson"`
				} `json:"cars"`
			} `json:"list"`
		} `json:"tp"`
	}

	if err := json.Unmarshal(response, &apiResponse); err != nil {
		return domain.TripsReturn{}, err
	}

	if apiResponse.Result != "OK" {
		return domain.TripsReturn{}, fmt.Errorf("API returned result: %s", apiResponse.Result)
	}

	var forwardTrips []domain.Trip
	var backTrips []domain.Trip

	for i, tp := range apiResponse.Tp {
		for _, tripData := range tp.List {
			depTime, err := time.Parse("02.01.2006 15:04", fmt.Sprintf("%s %s", tripData.Date0, tripData.Time0))
			if err != nil {
				return domain.TripsReturn{}, err
			}
			arrTime, err := time.Parse("02.01.2006 15:04", fmt.Sprintf("%s %s", tripData.Date1, tripData.Time1))
			if err != nil {
				return domain.TripsReturn{}, err
			}

			trip := domain.Trip{
				ID:          uuid.New().String(),
				RouteID:     tripData.Route0,
				TrainNumber: tripData.Number,
				Schedule: domain.Schedule{
					DepartureTime: depTime,
					ArrivalTime:   arrTime,
				},
				DepartureTime: depTime,
				ArrivalTime:   arrTime,
				TrainType:     tripData.Brand,
				Status:        "Запланирован",
			}

			if i == 0 {
				forwardTrips = append(forwardTrips, trip)
			} else if i == 1 {
				backTrips = append(backTrips, trip)
			}
		}
	}

	return domain.TripsReturn{
		Forward: forwardTrips,
		Back:    backTrips,
	}, nil
}

func (api *RzdAPI) TrainCarriages(ctx context.Context, params domain.TrainCarriagesParams) ([]domain.Carriage, error) {
	// Реализация получения вагонов выбранного поезда

	requestParams := map[string]interface{}{
		"layer_id": 5764,
		"dir":      params.Dir,
		"code0":    params.Code0,
		"code1":    params.Code1,
		"dt0":      params.Dt0.Format("02.01.2006"),
		"time0":    params.Time0.Format("15:04"),
		"tnum0":    params.Tnum0,
	}

	response, err := api.client.Post(ctx, "https://pass.rzd.ru/timetable/public/"+api.config.Language, requestParams)
	if err != nil {
		return nil, err
	}

	var apiResponse struct {
		Result string `json:"result"`
		Lst    []struct {
			Cars []struct {
				Cnumber    string  `json:"cnumber"`
				Type       string  `json:"type"`
				TypeLoc    string  `json:"typeLoc"`
				ClsType    string  `json:"clsType"`
				Tariff     float64 `json:"tariff"`
				TariffServ float64 `json:"tariffServ"`
				Seats      []struct {
					Places []string `json:"places"`
					Tariff float64  `json:"tariff"`
					Type   string   `json:"type"`
					Free   int      `json:"free"`
					Label  string   `json:"label"`
				} `json:"seats"`
			} `json:"cars"`
			FunctionBlocks interface{} `json:"functionBlocks"`
		} `json:"lst"`
		Schemes          interface{} `json:"schemes"`
		InsuranceCompany interface{} `json:"insuranceCompany"`
	}

	if err := json.Unmarshal(response, &apiResponse); err != nil {
		return nil, err
	}

	if apiResponse.Result != "OK" {
		return nil, fmt.Errorf("API returned result: %s", apiResponse.Result)
	}

	var carriages []domain.Carriage
	for _, lst := range apiResponse.Lst {
		for _, car := range lst.Cars {
			c := domain.Carriage{
				Number:        car.Cnumber,
				Type:          car.Type,
				TypeFullName:  car.TypeLoc,
				ClassType:     car.ClsType,
				Tariff:        car.Tariff,
				ServiceTariff: car.TariffServ,
			}

			for _, seat := range car.Seats {
				s := domain.Seat{
					Number: seat.Type, // Возможно, необходимо уточнение
					Type:   seat.Type,
					Tariff: seat.Tariff,
					Label:  seat.Label,
					Free:   seat.Free > 0,
				}
				c.Seats = append(c.Seats, s)
			}

			carriages = append(carriages, c)
		}
	}

	return carriages, nil
}

func (api *RzdAPI) TrainStationList(ctx context.Context, params domain.TrainStationListParams) ([]domain.Station, error) {
	// Реализация получения списка станций для маршрута

	requestParams := map[string]interface{}{
		"STRUCTURE_ID": 704,
		"trainNumber":  params.TrainNumber,
		"depDate":      params.DepDate.Format("02.01.2006"),
	}

	response, err := api.client.Get(ctx, "https://pass.rzd.ru/ticket/services/route/basicRoute", requestParams)
	if err != nil {
		return nil, err
	}

	var apiResponse struct {
		Data struct {
			TrainInfo struct {
				Number string `json:"number"`
			} `json:"trainInfo"`
			Routes []struct {
				Station     string `json:"station"`
				Code        string `json:"code"`
				ArvTime     string `json:"arvTime"`
				WaitingTime string `json:"waitingTime"`
				DepTime     string `json:"depTime"`
				Distance    int    `json:"distance"`
			} `json:"routes"`
		} `json:"data"`
	}

	if err := json.Unmarshal(response, &apiResponse); err != nil {
		return nil, err
	}

	var stations []domain.Station
	for _, route := range apiResponse.Data.Routes {
		st := domain.Station{
			Name: route.Station,
			ID:   route.Code,
			// Дополнительно можно добавить время прибытия, отправления и дистанцию
		}
		stations = append(stations, st)
	}

	return stations, nil
}

func (api *RzdAPI) StationCode(ctx context.Context, params domain.StationCodeParams) ([]domain.Station, error) {
	// Реализация получения кодов станций

	requestParams := map[string]interface{}{
		"lang":            api.config.Language,
		"stationNamePart": params.StationNamePart,
		"compactMode":     params.CompactMode,
	}

	response, err := api.client.Get(ctx, "https://pass.rzd.ru/suggester", requestParams)
	if err != nil {
		return nil, err
	}

	var apiResponse []struct {
		N string `json:"n"`
		C string `json:"c"`
	}

	if err := json.Unmarshal(response, &apiResponse); err != nil {
		return nil, err
	}

	var stationCodes []domain.Station
	for _, station := range apiResponse {
		if len(station.N) >= len(params.StationNamePart) &&
			station.N[:len(params.StationNamePart)] == params.StationNamePart {
			sc := domain.Station{
				Name: station.N,
				ID:   station.C,
			}
			stationCodes = append(stationCodes, sc)
		}
	}

	return stationCodes, nil
}
