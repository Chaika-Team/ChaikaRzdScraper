// internal/service/interface.go
package service

import (
	"context"

	"github.com/Chaika-Team/rzd-api/internal/domain"
)

// Service интерфейс сервиса для работы с маршрутами поездов и вагонами поездов РЖД
type Service interface {
	GetTrainRoutes(ctx context.Context, params domain.GetTrainRoutesParams) ([]domain.TrainRoute, error)
	GetTrainRoutesReturn(ctx context.Context, params domain.GetTrainRoutesReturnParams) ([]domain.TrainRoute, []domain.TrainRoute, error)
	GetTrainCarriages(ctx context.Context, params domain.GetTrainCarriagesParams) (domain.TrainCarriagesResponse, error)
	GetTrainStationList(ctx context.Context, params domain.GetTrainStationListParams) (domain.TrainStationListResponse, error)
	GetStationCode(ctx context.Context, params domain.GetStationCodeParams) ([]domain.StationCode, error)
}
