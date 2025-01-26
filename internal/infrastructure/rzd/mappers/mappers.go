package mappers

import (
	"github.com/Chaika-Team/rzd-api/internal/domain"
	"github.com/Chaika-Team/rzd-api/internal/infrastructure/rzd/schemas"
	"strconv"
	"time"
)

// MapTrainRouteResponse маппит ответ API маршрутов в доменную модель
func MapTrainRouteResponse(response schemas.TrainRouteResponse) ([]domain.TrainRoute, error) {
	var routes []domain.TrainRoute

	for _, tp := range response.TP {
		for _, trainList := range tp.List {
			for _, cases := range trainList.Cases {
				for _, trainCase := range cases {
					duration, _ := time.ParseDuration(trainCase.TimeInWay + "m") // "12:30" -> 12h30m -> 12h30m0s Duration
					departure, _ := parseDateTime(trainCase.Date0, trainCase.Time0)
					arrival, _ := parseDateTime(trainCase.Date1, trainCase.Time1)

					route := domain.TrainRoute{
						TrainNumber: trainCase.Number,
						Duration:    duration,
						Brand:       trainCase.Brand,
						Carrier:     trainCase.Carrier,
						From: domain.Station{
							Name: trainCase.Station0,
							Code: trainCase.Code0,
						},
						To: domain.Station{
							Name: trainCase.Station1,
							Code: trainCase.Code1,
						},
						Departure: departure,
						Arrival:   arrival,
						Cars:      mapTrainCarriages(trainCase.Cars),
					}

					routes = append(routes, route)
				}
			}
		}
	}

	return routes, nil
}

// mapTrainCarriages маппит список вагонов
func mapTrainCarriages(carriages []schemas.Carriage) []domain.Carriage {
	var result []domain.Carriage

	for _, car := range carriages {
		carriage := domain.Carriage{
			Number:   strconv.Itoa(car.Itype),
			Type:     car.TypeLoc,
			Class:    car.ServCls,
			Tariff:   float64(car.Tariff),
			Disabled: car.DisabledPerson,
			Seats:    mapSeats(car),
		}

		result = append(result, carriage)
	}

	return result
}

// mapSeats маппит места вагона
func mapSeats(car schemas.Carriage) []domain.Seat {
	return []domain.Seat{
		{
			Places: []string{}, // API может не возвращать конкретные места, это можно уточнить.
			Tariff: float64(car.Tariff),
			Type:   car.TypeLoc,
			Free:   int32(car.FreeSeats),
			Label:  car.Type,
		},
	}
}

// mapTrainStationList маппит ответ списка станций
func mapTrainStationList(response schemas.TrainStationListResponse) domain.TrainStationListResponse {
	var routes []domain.RouteInfo

	for _, route := range response.Data.Routes {
		arrivalTime, _ := parseTime(route.ArvTime)
		departureTime, _ := parseTime(route.DepTime)

		routes = append(routes, domain.RouteInfo{
			Station: domain.Station{
				Name: route.Station,
			},
			ArrivalTime:   arrivalTime,
			DepartureTime: departureTime,
			Distance:      route.Distance,
		})
	}

	return domain.TrainStationListResponse{
		Train: domain.TrainInfo{
			Number: response.Data.TrainInfo.Number,
		},
		Routes: routes,
	}
}

// parseDateTime парсит дату и время из строки
func parseDateTime(dateStr, timeStr string) (time.Time, error) {
	layout := "02.01.2006 15:04"
	return time.Parse(layout, dateStr+" "+timeStr)
}

// parseTime парсит только время
func parseTime(timeStr string) (time.Time, error) {
	layout := "15:04"
	return time.Parse(layout, timeStr)
}
