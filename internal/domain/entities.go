// internal/domain/entities.go
package domain

// TrainRoute доменная модель маршрута поезда
type TrainRoute struct {
	Route0        string     // Название станции отправления
	Route1        string     // Название станции назначения
	DepartureDate string     // Дата отправления
	DepartureTime string     // Время отправления
	Number        string     // Номер поезда
	From          string     // Название станции отправления
	Where         string     // Название станции назначения
	ArrivalDate   string     // Дата прибытия
	FromCode      string     // Код станции отправления
	WhereCode     string     // Код станции назначения
	ArrivalTime   string     // Время прибытия
	TimeInWay     string     // Время в пути
	Brand         string     // Бренд поезда
	Carrier       string     // Перевозчик
	Cars          []Carriage // Cписок вагонов поезда
}

// Carriage доменная модель вагона поезда
type Carriage struct {
	CNumber string // Номер вагона
	Type    string // Тип вагона
	TypeLoc string
	ClsType string  // Класс вагона
	Tariff  float32 // Тариф
	Seats   []Seat  // Список мест в вагоне
}

type Seat struct {
	Places []string // Список мест
	Tariff float32  // Тариф
	Type   string   // Тип
	Free   int32    // Количество свободных мест
	Label  string
}

type TrainCarriagesResponse struct {
	Cars           []Carriage
	FunctionBlocks []string
	Schemes        []string
	Companies      []string
}

type TrainInfo struct {
	Number string
	// Добавить другие поля по необходимости
}

type RouteInfo struct {
	Station       string
	ArrivalTime   string
	DepartureTime string
	// Добавить другие поля по необходимости
}

type TrainStationListResponse struct {
	Train  TrainInfo
	Routes []RouteInfo
}

type StationCode struct {
	Station string
	Code    string
}

type GetTrainRoutesParams struct {
	Code0      string
	Code1      string
	Dir        int32
	Tfl        int32
	CheckSeats int32
	Dt0        string
	Md         int32
}

type GetTrainRoutesReturnParams struct {
	Code0      string
	Code1      string
	Dir        int32
	Tfl        int32
	CheckSeats int32
	Dt0        string
	Dt1        string
}

type GetTrainCarriagesParams struct {
	Code0 string
	Code1 string
	Tnum0 string
	Time0 string
	Dt0   string
	Dir   int32
}

type GetTrainStationListParams struct {
	TrainNumber string
	DepDate     string
}

type GetStationCodeParams struct {
	StationNamePart string
	CompactMode     string
}
