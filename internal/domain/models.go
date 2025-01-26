// internal/domain/models.go
package domain

import (
	"time"
)

// TrainRoute представляет маршрут поезда
type TrainRoute struct {
	TrainNumber string        // Номер поезда
	Duration    time.Duration // Время в пути
	Brand       string        // Бренд поезда
	Carrier     string        // Перевозчик

	From      Station // Станция отправления
	To        Station // Станция прибытия
	Departure time.Time
	Arrival   time.Time

	Cars []Carriage // Список вагонов поезда
}

// Station представляет железнодорожную станцию
type Station struct {
	Name string // Название станции
	Code int    // Код станции
}

// Carriage представляет вагон поезда
type Carriage struct {
	Number    string // Номер вагона
	Type      string // Тип вагона (купе, плацкарт и т.д.)
	TypeLabel string // Полное наименование типа вагона
	Class     string // Класс вагона (2Л, 2Э и т.д.)
	Tariff    int    // Стоимость билета
	TariffEx  int    // Тариф за место 2 TODO: что это
	FreeSeats int    // Количество свободных мест
	Seats     []Seat // Список мест в вагоне
	Disabled  bool   // Места для инвалидов
}

// Seat представляет место в вагоне //TODO не готово
type Seat struct {
	Places []string // Список свободных мест
	Tariff int      // Тариф за место
	Type   string   // Тип места (верхние, нижние и т.д.)
	Label  string   // Полное наименование места
}

// TrainCarriagesResponse представляет ответ на запрос вагонов
type TrainCarriagesResponse struct {
	Cars           []Carriage // Список вагонов
	FunctionBlocks []string   // Функциональные блоки
	Schemes        []string   // Схемы вагонов
	Companies      []string   // Компании перевозчики
}

// TrainInfo представляет информацию о поезде
type TrainInfo struct {
	Number string // Номер поезда
	// Добавить другие поля по необходимости
}

// RouteInfo представляет информацию о маршруте
type RouteInfo struct {
	Station       Station   // Станция
	ArrivalTime   time.Time // Время прибытия
	DepartureTime time.Time // Время отправления
	Distance      int       // Пройденная дистанция
}

// TrainStationListResponse представляет ответ на запрос списка станций
type TrainStationListResponse struct {
	Train  TrainInfo   // Информация о поезде
	Routes []RouteInfo // Список станций
}

// StationCode представляет код станции
type StationCode struct {
	Station string // Название станции
	Code    int    // Код станции
}

// GetTrainRoutesParams представляет параметры для запроса маршрутов
type GetTrainRoutesParams struct {
	FromCode   int       // Код станции отправления
	ToCode     int       // Код станции прибытия
	Direction  Direction // Направление
	TrainType  TrainType // Тип поезда
	CheckSeats bool      // Проверка наличия мест
	FromDate   time.Time // Дата отправления
	WithChange bool      // С пересадками
}

// GetTrainRoutesReturnParams представляет параметры для запроса маршрутов туда-обратно
type GetTrainRoutesReturnParams struct {
	GetTrainRoutesParams
	ToDate time.Time // Дата возвращения
}

// GetTrainCarriagesParams представляет параметры для запроса вагонов
type GetTrainCarriagesParams struct {
	TrainNumber string    // Номер поезда
	Direction   Direction // Направление
	FromCode    int       // Код станции отправления
	FromTime    time.Time // Время отправления
	FromDate    time.Time // Дата отправления
	ToCode      int       // Код станции прибытия
}

// GetTrainStationListParams представляет параметры для запроса списка станций
type GetTrainStationListParams struct {
	TrainNumber string    // Номер поезда
	FromDate    time.Time // Дата отправления
}

// GetStationCodeParams представляет параметры для запроса кодов станций
type GetStationCodeParams struct {
	StationNamePart string // Часть названия станции
	CompactMode     bool   // Компактный режим
}
