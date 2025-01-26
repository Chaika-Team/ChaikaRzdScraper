package mappers

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Chaika-Team/rzd-api/internal/domain"
	"github.com/Chaika-Team/rzd-api/internal/infrastructure/rzd/schemas"
)

// MapTrainRouteResponse маппит ответ API маршрутов в доменную модель
func MapTrainRouteResponse(response schemas.TrainRouteResponse) ([]domain.TrainRoute, error) {
	var routes []domain.TrainRoute

	// Перебор всех TP в ответе
	for _, tp := range response.TP {
		// Перебор всех поездов в списке
		for _, train := range tp.List {
			// Парсинг времени в пути
			duration, err := parseDuration(train.TimeInWay)
			if err != nil {
				return nil, fmt.Errorf("failed to parse timeInWay: %v", err)
			}

			// Парсинг времени отправления и прибытия
			departure, err := parseDateTime(train.Date0, train.Time0)
			if err != nil {
				return nil, fmt.Errorf("failed to parse departure time: %v", err)
			}

			arrival, err := parseDateTime(train.Date1, train.Time1)
			if err != nil {
				return nil, fmt.Errorf("failed to parse arrival time: %v", err)
			}

			// Маппинг маршрута
			route := domain.TrainRoute{
				TrainNumber: train.Number,
				Duration:    duration,
				Brand:       train.Brand,
				Carrier:     train.Carrier,
				From: domain.Station{
					Name: train.Station0,
					Code: train.Code0,
				},
				To: domain.Station{
					Name: train.Station1,
					Code: train.Code1,
				},
				Departure: departure,
				Arrival:   arrival,
				Cars:      mapTrainCarriages(train.Cars),
			}

			// Обработка seatCars (если они есть)
			if len(train.SeatCars) > 0 {
				seatCarriages := mapTrainSeatCarriages(train.SeatCars)
				route.Cars = append(route.Cars, seatCarriages...)
			}

			routes = append(routes, route)
		}
	}

	return routes, nil
}

// mapTrainCarriages маппит список вагонов
func mapTrainCarriages(carriages []schemas.Carriage) []domain.Carriage {
	var result []domain.Carriage

	for _, car := range carriages {
		carriage := domain.Carriage{
			Number:    strconv.Itoa(car.Itype),
			Type:      car.Type,
			TypeLabel: car.TypeLoc,
			Class:     car.ServCls,
			Tariff:    car.Tariff,
			Disabled:  car.DisabledPerson,
			Seats:     []domain.Seat{},
		}

		result = append(result, carriage)
	}

	return result
}

func mapTrainSeatCarriages(carriages []schemas.SeatCarriage) []domain.Carriage {
	var result []domain.Carriage

	for _, car := range carriages {

		tariff, err := strconv.Atoi(car.Tariff)
		if err != nil {
			log.Printf("failed to parse tariff: %v", err)
			tariff = 0
		}
		tariff2, err := strconv.Atoi(car.Tariff2)
		if err != nil {
			log.Printf("failed to parse tariff2: %v", err)
			tariff2 = 0
		}

		carriage := domain.Carriage{
			Number:    strconv.Itoa(car.Itype),
			Type:      car.Type,
			TypeLabel: car.TypeLoc,
			Class:     car.ServCls,
			Tariff:    tariff,
			TariffEx:  tariff2,
			Disabled:  false, // Как понимаю SeatCarriage сделано для бизнес-класса, где нет мест для инвалидов
			FreeSeats: car.FreeSeats,
			Seats:     []domain.Seat{},
		}

		result = append(result, carriage)
	}

	return result
}

// mapSeats маппит места вагона // TODO не готово
func mapSeats(car schemas.Carriage) []domain.Seat {

	return []domain.Seat{
		{
			Places: []string{}, // API может не возвращать конкретные места, это можно уточнить.
			Tariff: car.Tariff,
			Type:   car.TypeLoc,
			Label:  car.Type,
		},
	}
}

// mapTrainStationList маппит ответ списка станций  // TODO не готово
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

// parseDuration преобразует строку времени в пути формата HH:mm в time.Duration
func parseDuration(durationStr string) (time.Duration, error) {
	parts := strings.Split(durationStr, ":")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid duration format: %s", durationStr)
	}
	hours, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("invalid hours in duration: %s", parts[0])
	}
	minutes, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, fmt.Errorf("invalid minutes in duration: %s", parts[1])
	}
	return time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute, nil
}

// parseTime парсит только время
func parseTime(timeStr string) (time.Time, error) {
	layout := "15:04"
	return time.Parse(layout, timeStr)
}
