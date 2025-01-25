// internal/usecase/station_usecase.go
package usecase

import (
	"context"

	"github.com/Chaika-Team/rzd-api/internal/domain"
	"github.com/Chaika-Team/rzd-api/internal/infrastructure/rzd"
)

type StationUseCase struct {
	RzdClient *rzd.RzdClient
}

func NewStationUseCase(client *rzd.RzdClient) *StationUseCase {
	return &StationUseCase{
		RzdClient: client,
	}
}

func (uc *StationUseCase) GetTrainStationList(ctx context.Context, params domain.GetTrainStationListParams) (domain.TrainStationListResponse, error) {
	return uc.RzdClient.GetTrainStationList(params)
}

func (uc *StationUseCase) GetStationCode(ctx context.Context, params domain.GetStationCodeParams) ([]domain.StationCode, error) {
	return uc.RzdClient.GetStationCode(params)
}
