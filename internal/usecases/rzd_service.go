// internal/usecases/rzd_service.go

package usecases

import (
	"context"

	"github.com/Chaika-Team/rzd-api/internal/domain"
)

type RzdService struct {
	api domain.RzdAPI
}

func NewRzdService(api domain.RzdAPI) *RzdService {
	return &RzdService{
		api: api,
	}
}

func (s *RzdService) GetTrainRoutes(ctx context.Context, params domain.TrainRoutesParams) ([]domain.Trip, error) {
	return s.api.TrainRoutes(ctx, params)
}

func (s *RzdService) GetTrainRoutesReturn(ctx context.Context, params domain.TrainRoutesReturnParams) (domain.TripsReturn, error) {
	return s.api.TrainRoutesReturn(ctx, params)
}

func (s *RzdService) GetTrainCarriages(ctx context.Context, params domain.TrainCarriagesParams) ([]domain.Carriage, error) {
	return s.api.TrainCarriages(ctx, params)
}

func (s *RzdService) GetTrainStationList(ctx context.Context, params domain.TrainStationListParams) ([]domain.Station, error) {
	return s.api.TrainStationList(ctx, params)
}

func (s *RzdService) GetStationCodes(ctx context.Context, params domain.StationCodeParams) ([]domain.Station, error) {
	return s.api.StationCode(ctx, params)
}
