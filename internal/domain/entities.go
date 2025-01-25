// internal/domain/entities.go
package domain

// TrainRoute доменная модель маршрута поезда
type TrainRoute struct {
	TrainNumber string // Номер поезда
	TimeInWay   string // Время в пути
	TrainBrand  string // Бренд поезда
	Carrier     string // Перевозчик

	From      string // Название станции отправления (САНКТ-ПЕТЕРБУРГ)
	FromShort string // Кодовое название станции отправления (С-ПЕТ-ЛАД)
	FromCode  string // Код станции отправления (2004000)
	FromDate  string // Дата отправления
	FromTime  string // Время отправления

	Where      string // Название станции прибытия (КИРОВ ПАСС)
	WhereShort string // Кодовое название станции назначения (ТЮМЕНЬ)
	WhereCode  string // Код станции прибытия (2060600)
	WhereDate  string // Дата прибытия
	WhereTime  string // Время прибытия

	Cars []Carriage // Список вагонов поезда
}

// Carriage доменная модель вагона поезда
type Carriage struct {
	CNumber   string  // Номер вагона
	Type      string  // Тип вагона
	TypeLoc   string  // Полное наименование (Плацкартный, СВ, Купе, Люкс)
	ClassType string  // Класс вагона (2Л, 2Э)
	Tariff    float32 // Стоимость билета
	Seats     []Seat  // Список мест в вагоне
}

type Seat struct {
	Places []string // Список свободных мест
	Tariff float32  // Тариф за место
	Type   string   // Сокращенное наименование места (up)
	Free   int32    // Количество свободных мест
	Label  string   // Полное наименование места (Верхние)
}

type TrainCarriagesResponse struct {
	Cars           []Carriage // Список вагонов
	FunctionBlocks []string   // Функциональные блоки
	Schemes        []string   // Схемы вагонов
	Companies      []string   // Компании перевозчики
}

type TrainInfo struct {
	Number string // Номер поезда
	// Добавить другие поля по необходимости
}

type RouteInfo struct {
	Station   string // Название станции
	WhereTime string // Время прибытия
	FromTime  string // Время отправления
	// Добавить другие поля по необходимости
}

type TrainStationListResponse struct {
	Train  TrainInfo   // Информация о поезде
	Routes []RouteInfo // Список станций
}

type StationCode struct {
	Station string // Название станции
	Code    string // Код станции
}

type GetTrainRoutesParams struct {
	FromCode   string // Код станции отправления
	WhereCode  string // Код станции прибытия
	Direction  int32  // Направление (0 - в одну сторону, 1 - туда и обратно)
	TrainType  int32  // Тип поезда (3 - поезда и электрички, 1 - только поезда, 2 - только электрички)
	CheckSeats int32  // Наличие свободных мест (0 - все, 1 - только с местами)
	FromDate   string // Дата отправления
	WithChange int32  // С пересадками (0 - без пересадок, 1 - с пересадками)
}

type GetTrainRoutesReturnParams struct {
	Direction  int32 // Направление (0 - в одну сторону, 1 - туда и обратно)
	TrainType  int32 // Тип поезда (3 - поезда и электрички, 1 - только поезда, 2 - только электрички)
	CheckSeats int32 // Наличие свободных мест (0 - все, 1 - только с местами)

	FromCode string // Код станции отправления
	FromDate string // Дата отправления

	WhereCode string // Код станции прибытия
	WhereDate string // Дата прибытия
}

type GetTrainCarriagesParams struct {
	TrainNumber string // Номер поезда (072Е)
	Direction   int32  // Направление (0 - в одну сторону, 1 - туда и обратно)
	FromCode    string // Код станции отправления (2004000)
	FromTime    string // Время отправления (00:00)
	FromDate    string // Дата отправления (28.12.2020)
	WhereCode   string // Код станции прибытия (2060600)
}

type GetTrainStationListParams struct {
	TrainNumber string // Номер поезда (072Е)
	FromDate    string // Дата отправления (28.12.2020)
}

type GetStationCodeParams struct {
	StationNamePart string // Часть названия станции, минимум 2 символа
	CompactMode     string // Компактный режим (по умолчанию 'y')
}
