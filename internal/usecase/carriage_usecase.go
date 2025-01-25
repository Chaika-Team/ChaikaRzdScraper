// internal/usecase/carriage_usecase.go
package usecase

import (
	"context"

	"github.com/Chaika-Team/rzd-api/internal/domain"
	"github.com/Chaika-Team/rzd-api/internal/infrastructure/rzd"
)

type CarriageUseCase struct {
	RzdClient *rzd.RzdClient
}

func NewCarriageUseCase(client *rzd.RzdClient) *CarriageUseCase {
	return &CarriageUseCase{
		RzdClient: client,
	}
}

func (uc *CarriageUseCase) GetTrainCarriages(ctx context.Context, params domain.GetTrainCarriagesParams) (domain.TrainCarriagesResponse, error) {
	carriages, err := uc.RzdClient.GetTrainCarriages(params)
	if err != nil {
		return domain.TrainCarriagesResponse{}, err
	}
	return carriages, nil
}

type GetTrainCarriagesParams struct {
	Code0 string
	Code1 string
	Tnum0 string
	Time0 string
	Dt0   string
	Dir   int32
}
