package mappers

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Chaika-Team/ChaikaRzdScraper/internal/domain"
	"github.com/Chaika-Team/ChaikaRzdScraper/internal/infrastructure/rzd/schemas"
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
				Carrier: domain.Carrier{
					Name: train.Carrier,
				},
				From: domain.Station{
					Name:      train.Station0,
					RouteName: train.Route0,
					Code:      train.Code0,
				},
				To: domain.Station{
					Name:      train.Station1,
					RouteName: train.Route1,
					Code:      train.Code1,
				},
				Departure: departure,
				Arrival:   arrival,
				CarTypes:  mapTrainCarriages(train.Cars),
			}

			// Обработка seatCars (если они есть)
			if len(train.SeatCars) > 0 {
				seatCarriages := mapTrainSeatCarriages(train.SeatCars)
				route.CarTypes = append(route.CarTypes, seatCarriages...)
			}

			routes = append(routes, route)
		}
	}

	return routes, nil
}

// mapTrainCarriages маппит список вагонов
func mapTrainCarriages(carriages []schemas.CarriageType) []domain.CarriageType {
	var result []domain.CarriageType

	for _, car := range carriages {
		carriage := domain.CarriageType{
			Type:           domain.CarSeatType(car.Itype),
			TypeShortLabel: car.Type,
			TypeLabel:      car.TypeLoc,
			Class:          car.ServCls,
			Tariff:         car.Tariff,
			Disabled:       car.DisabledPerson,
		}

		result = append(result, carriage)
	}

	return result
}

func mapTrainSeatCarriages(carriages []schemas.SeatCarriageType) []domain.CarriageType {
	var result []domain.CarriageType

	for _, car := range carriages {

		tariff, err := strconv.Atoi(car.Tariff)
		if err != nil {
			log.Printf("failed to parse tariff for car type %s (defaulting to 0): %v", car.Type, err)
			tariff = 0
		}
		tariff2 := 0
		if car.Tariff2 != "" {
			tariff2, err = strconv.Atoi(car.Tariff2)
			if err != nil {
				log.Printf("failed to parse tariff2 for car type %s (defaulting to 0): %v", car.Type, err)
				tariff2 = 0
			}
		}

		carriage := domain.CarriageType{
			Type:        domain.CarSeatType(car.Itype),
			TypeLabel:   car.TypeLoc,
			Class:       car.ServCls,
			Tariff:      tariff,
			TariffExtra: tariff2,
			Disabled:    false, // Как понимаю SeatCarriageType сделано для бизнес-класса, где нет мест для инвалидов
			FreeSeats:   car.FreeSeats,
		}

		result = append(result, carriage)
	}

	return result
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
