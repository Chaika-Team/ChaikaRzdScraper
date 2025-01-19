// internal/domain/interfaces.go

package domain

import (
	"context"
	"time"
)

// RzdAPI интерфейс для взаимодействия с RZD API.
type RzdAPI interface {
	TrainRoutes(ctx context.Context, params TrainRoutesParams) ([]Trip, error)
	TrainRoutesReturn(ctx context.Context, params TrainRoutesReturnParams) (TripsReturn, error)
	TrainCarriages(ctx context.Context, params TrainCarriagesParams) ([]Carriage, error)
	TrainStationList(ctx context.Context, params TrainStationListParams) ([]Station, error)
	StationCode(ctx context.Context, params StationCodeParams) ([]Station, error)
}

// TrainRoutesParams Используемые параметры для методов
type TrainRoutesParams struct {
	Dir        int
	Tfl        int
	CheckSeats int
	Code0      string
	Code1      string
	Dt0        time.Time
	Md         int
}

type TrainRoutesReturnParams struct {
	Dir        int
	Tfl        int
	CheckSeats int
	Code0      string
	Code1      string
	Dt0        time.Time
	Dt1        time.Time
}

type TrainCarriagesParams struct {
	Dir   int
	Code0 string
	Code1 string
	Dt0   time.Time
	Time0 time.Time
	Tnum0 string
}

type TrainStationListParams struct {
	TrainNumber string
	DepDate     time.Time
}

type StationCodeParams struct {
	StationNamePart string
	CompactMode     string
}

// TripsReturn представляет маршруты туда-обратно.
type TripsReturn struct {
	Forward []Trip
	Back    []Trip
}
