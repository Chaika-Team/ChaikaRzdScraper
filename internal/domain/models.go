// internal/domain/models.go
package domain

import (
	"time"
)

// TrainRoute представляет маршрут поезда, полученный при запросе рейсов.
// Содержит общую информацию о поезде и список типов вагонов (агрегированные данные).
type TrainRoute struct {
	TrainNumber string        // Номер поезда
	TrainType   TrainType     // Тип поезда (например, поезд или электричка)
	Duration    time.Duration // Время в пути
	Brand       string        // Бренд поезда
	Carrier     Carrier       // Перевозчик

	From      Station   // Станция отправления
	To        Station   // Станция прибытия
	Departure time.Time // Время отправления (дата и время)
	Arrival   time.Time // Время прибытия (дата и время)

	CarTypes []CarriageType // Список типов вагонов поезда (агрегированные данные)

	Cars []Car // Список конкретных вагонов поезда
}

// Station представляет железнодорожную станцию.
type Station struct {
	Name      string // Полное название станции, например "САНКТ-ПЕТЕРБУРГ-ГЛАВН. (МОСКОВСКИЙ ВОКЗАЛ)"
	RouteName string // Краткое название станции в маршруте, например "С-ПЕТЕР-ГЛ"
	Code      int    // Код станции (например, 2004000, 2000000)
	Level     int    // Уровень станции (0-5, 5 - самый высокий)
	Score     int    // Значение сортировки (0-5, 5 - самое высокое)
}

// CarriageType представляет агрегированные данные о типах вагонов, полученные из запроса маршрутов.
// Используется для отображения общего состава поезда (без привязки к конкретным вагонам).
type CarriageType struct {
	Type           CarSeatType // Тип посадочных мест в вагоне (например, плацкарт, купе, люкс)
	TypeShortLabel string      // Краткое наименование типа вагона
	TypeLabel      string      // Полное наименование типа вагона
	Class          string      // Класс вагона (например, "2Л", "2Э")
	Tariff         int         // Стоимость билета для данного типа вагона
	TariffExtra    int         // Дополнительный тариф (если имеется)
	FreeSeats      int         // Общее количество свободных мест в вагонах этого типа
	Disabled       bool        // Флаг: есть ли специальные места для инвалидов
}

// Car представляет один конкретный вагон поезда с детальной информацией.
type Car struct {
	CarNumber          string        // Номер вагона
	Type               string        // Тип вагона (например, "Купе", "Плац", "Люкс")
	CategoryLabelLocal string        // Категория вагона (например, "Купе")
	TypeLabel          string        // Полное наименование типа вагона, например "Купе"
	CategoryCode       string        // Код категории вагона
	CarTypeID          int           // Идентификатор категории вагона
	CarType            int           // Тип вагона, обычно также код
	Letter             string        // Буква вагона
	ClassType          string        // Тип класса вагона (например, "2Ш")
	Services           []Service     // Список услуг, предоставляемых в вагоне
	Tariff             int           // Стоимость билета в вагоне
	Tariff2            int           // Дополнительный тариф (если имеется)
	Carrier            Carrier       // Перевозчик
	CarNumeration      CarNumeration // Нумерация вагона // TODO почему это в вагоне а не в поезде?

}

type Service struct {
	ID          string // Идентификатор услуги
	Name        string // Название услуги (с иконкой)
	Description string // Описание услуги
}

type Carrier struct {
	ID   string // Идентификатор перевозчика
	Name string // Название перевозчика
}

// GetTrainRoutesParams представляет параметры для запроса маршрутов
type GetTrainRoutesParams struct {
	FromCode   int             // Код станции отправления
	ToCode     int             // Код станции прибытия
	Direction  Direction       // Направление
	TrainType  TrainSearchType // Тип поезда
	CheckSeats bool            // Проверка наличия мест
	FromDate   time.Time       // Дата отправления
	WithChange bool            // С пересадками
}

// GetTrainCarriagesParams представляет параметры для запроса вагонов
type GetTrainCarriagesParams struct {
	TrainNumber string    // Номер поезда
	Direction   Direction // Направление
	FromCode    int       // Код станции отправления
	FromTime    time.Time // Время отправления
	ToCode      int       // Код станции прибытия
}

// SearchStationParams представляет параметры для поиска станций по части названия.
type SearchStationParams struct {
	Query       string // Поисковый запрос, например "ЧЕБ"
	CompactMode bool   // Флаг компактного режима
}
