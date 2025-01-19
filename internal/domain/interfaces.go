// internal/domain/interfaces.go
package domain

import (
	"context"
	"time"
)

// RzdAPI задаёт методы в стиле старого PHP: возвращаем JSON-строку
type RzdAPI interface {
	TrainRoutes(ctx context.Context, params TrainRoutesParams) (string, error)
	TrainRoutesReturn(ctx context.Context, params TrainRoutesReturnParams) (string, error)
	TrainCarriages(ctx context.Context, params TrainCarriagesParams) (string, error)
	TrainStationList(ctx context.Context, params TrainStationListParams) (string, error)
	StationCode(ctx context.Context, params StationCodeParams) (string, error)
}

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
