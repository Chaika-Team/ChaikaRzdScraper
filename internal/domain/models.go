// internal/domain/models.go
package domain

import "time"

// Набор струтктур, если нужно где-то использовать (необязательно):
type Route struct {
	ID                     string
	CodeName               string
	StartStationID         string
	EndStationID           string
	IntermediateStationIDs []string
	Length                 float64
	TravelTime             time.Duration
	RouteType              string
	Status                 string
}

// Trip — пример доменной модели, если вдруг понадобится
type Trip struct {
	ID                    string
	RouteID               string
	TrainNumber           string
	Schedule              Schedule
	DepartureTime         time.Time
	ArrivalTime           time.Time
	TrainType             string
	CarriageConfiguration CarriageConfiguration
	Status                string
}

type Schedule struct {
	DepartureTime time.Time
	ArrivalTime   time.Time
}

type CarriageConfiguration struct {
	Type  string
	Seats int
}

type Station struct {
	ID             string
	Name           string
	Latitude       float64
	Longitude      float64
	Infrastructure string
	Type           string
}

type Carriage struct {
	Number        string
	Type          string
	TypeFullName  string
	ClassType     string
	Tariff        float64
	ServiceTariff float64
	Seats         []Seat
}

type Seat struct {
	Number string
	Type   string
	Tariff float64
	Label  string
	Free   bool
}

// Если вдруг потребуется
type Profile struct {
	ID       string
	Username string
	Email    string
	FullName string
}

// Для "туда-обратно"
type TripsReturn struct {
	Forward interface{} `json:"forward"`
	Back    interface{} `json:"back"`
}
