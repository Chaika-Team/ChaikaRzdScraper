// internal/domain/models.go

package domain

import "time"

// Route представляет маршрут поезда.
type Route struct {
	ID                     string
	CodeName               string
	StartStationID         string
	EndStationID           string
	IntermediateStationIDs []string
	Length                 float64 // в км
	TravelTime             time.Duration
	RouteType              string // Например, пассажирский, экспресс
	Status                 string // Активен, приостановлен, закрыт
}

// Trip представляет рейс поезда.
type Trip struct {
	ID                    string
	RouteID               string
	TrainNumber           string
	Schedule              Schedule
	DepartureTime         time.Time
	ArrivalTime           time.Time
	TrainType             string
	CarriageConfiguration CarriageConfiguration
	Status                string // Запланирован, отменен, задержан, выполнен
}

// Schedule представляет расписание рейса.
type Schedule struct {
	DepartureTime time.Time
	ArrivalTime   time.Time
}

// CarriageConfiguration представляет конфигурацию вагонов.
type CarriageConfiguration struct {
	Type  string
	Seats int
}

// Station представляет железнодорожную станцию.
type Station struct {
	ID             string
	Name           string
	Latitude       float64
	Longitude      float64
	Infrastructure string
	Type           string // Основная, промежуточная, конечная
}

// Carriage представляет вагон поезда.
type Carriage struct {
	Number        string
	Type          string
	TypeFullName  string
	ClassType     string
	Tariff        float64
	ServiceTariff float64
	Seats         []Seat
}

// Seat представляет место в вагоне.
type Seat struct {
	Number string
	Type   string
	Tariff float64
	Label  string
	Free   bool
}

// Profile представляет данные пользователя.
type Profile struct {
	ID       string
	Username string
	Email    string
	FullName string
	// Другие поля по необходимости
}
