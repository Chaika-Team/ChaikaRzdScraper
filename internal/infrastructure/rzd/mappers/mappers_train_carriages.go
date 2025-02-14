// internal/infrastructure/rzd/mappers/train_carriages_mapper.go
package mappers

import (
	"fmt"
	"strconv"

	"github.com/Chaika-Team/rzd-api/internal/domain"
	"github.com/Chaika-Team/rzd-api/internal/infrastructure/rzd/schemas"
)

// MapTrainCarriagesResponse преобразует схему TrainCarriagesResponse в срез доменных моделей Car.
// Игнорируются поля: Timestamp, FunctionBlocks, а также некоторые дополнительные флаги, не входящие в доменную модель.
func MapTrainCarriagesResponse(resp schemas.TrainCarriagesResponse) ([]domain.Car, error) {
	var cars []domain.Car

	// Проверяем, что список результатов не пуст.
	if len(resp.Lst) == 0 {
		return nil, fmt.Errorf("response contains no train results")
	}

	// Обычно в ответе возвращается один поезд, но мы пройдёмся по всем найденным
	for _, trainResult := range resp.Lst {
		// Проходим по каждому вагону в поезде
		for _, carSchema := range trainResult.Cars {
			// Маппим список услуг
			var serviceList []domain.Service
			for _, s := range carSchema.Services {
				// Преобразуем числовой идентификатор в строку
				serviceList = append(serviceList, domain.Service{
					ID:          strconv.Itoa(s.ID), // Преобразование int -> string
					Name:        s.Name,
					Description: s.Description,
					// Поле HasImage мы не используем в доменной модели (оно игнорируется)
				})
			}

			// Преобразуем тарифы (из строки в int)
			tariff, err := strconv.Atoi(carSchema.Tariff)
			if err != nil {
				// Если конвертация не удалась, тариф считаем равным 0
				tariff = 0
			}
			tariff2, err := strconv.Atoi(carSchema.Tariff2)
			if err != nil {
				tariff2 = 0
			}

			// Маппинг перевозчика: в доменной модели Carrier имеет поля ID и Name,
			// где ID – строка. В схеме перевозчика передаётся как "carrier" (имя) и "carrierId" (число).
			carrier := domain.Carrier{
				ID:   strconv.Itoa(carSchema.CarrierId),
				Name: carSchema.Carrier,
			}

			// Маппинг нумерации вагона: поле CarNumeration может быть nil, если отсутствует.
			var carNumeration = mapCarNumeration(carSchema.CarNumeration)

			// Маппинг списка мест в вагоне, просто посчитаем их количество
			var freeSeats = len(carSchema.Seats)

			// Собираем данные о конкретном вагоне в доменную модель.
			// Поля, которые не используются в доменной модели (например, AddSigns, IntServiceClass и др.) игнорируются.
			car := domain.Car{
				CarNumber:          carSchema.Cnumber,
				Type:               carSchema.Type,
				CategoryLabelLocal: carSchema.CatLabelLoc,
				TypeLabel:          carSchema.TypeLoc,
				CategoryCode:       carSchema.CatCode,
				CarTypeID:          carSchema.Ctypei,
				CarType:            carSchema.Ctype,
				Letter:             carSchema.Letter,
				ClassType:          carSchema.ClsType,
				Services:           serviceList,
				Tariff:             tariff,
				Tariff2:            tariff2,
				Carrier:            carrier,
				CarNumeration:      carNumeration,
				FreeSeats:          freeSeats,
			}

			cars = append(cars, car)
		}
	}

	return cars, nil
}

// mapCarNumeration преобразует строковое представление нумерации вагона в CarNumeration.
func mapCarNumeration(value *string) domain.CarNumeration {
	if value == nil {
		return domain.Unknown
	}

	switch *value {
	case "FromHead":
		return domain.Head
	case "FromTail":
		return domain.Tail
	default:
		return domain.Unknown
	}
}
