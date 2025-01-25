// internal/usecase/route_usecase.go
package usecase

import (
	"context"

	"github.com/Chaika-Team/rzd-api/internal/domain"
	"github.com/Chaika-Team/rzd-api/internal/infrastructure/rzd"
)

type RouteUseCase struct {
	RzdClient *rzd.RzdClient
}

func NewRouteUseCase(client *rzd.RzdClient) *RouteUseCase {
	return &RouteUseCase{
		RzdClient: client,
	}
}

func (uc *RouteUseCase) GetTrainRoutes(ctx context.Context, params domain.GetTrainRoutesParams) ([]domain.TrainRoute, error) {
	return uc.RzdClient.GetTrainRoutes(params)
}

func (uc *RouteUseCase) GetTrainRoutesReturn(ctx context.Context, params domain.GetTrainRoutesReturnParams) ([]domain.TrainRoute, []domain.TrainRoute, error) {
	return uc.RzdClient.GetTrainRoutesReturn(params)
}
