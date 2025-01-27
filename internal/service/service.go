// internal/service/service.go
package service

import (
	"context"

	"github.com/Chaika-Team/rzd-api/internal/domain"
	"github.com/Chaika-Team/rzd-api/internal/infrastructure/rzd"
)

// mainService реализует интерфейс Service
type mainService struct {
	rzdClient *rzd.Client
}

// New возвращает новый экземпляр сервиса
func New(rzdClient *rzd.Client) Service {
	return &mainService{rzdClient: rzdClient}
}

// GetTrainRoutes реализует метод получения маршрутов поездов
func (s *mainService) GetTrainRoutes(ctx context.Context, params domain.GetTrainRoutesParams) ([]domain.TrainRoute, error) {
	return s.rzdClient.GetTrainRoutes(ctx, params)
}

// GetTrainRoutesReturn реализует метод получения маршрутов туда и обратно
func (s *mainService) GetTrainRoutesReturn(ctx context.Context, params domain.GetTrainRoutesReturnParams) ([]domain.TrainRoute, []domain.TrainRoute, error) {
	return s.rzdClient.GetTrainRoutesReturn(ctx, params)
}

// GetTrainCarriages реализует метод получения информации о вагонах
func (s *mainService) GetTrainCarriages(ctx context.Context, params domain.GetTrainCarriagesParams) (domain.TrainCarriagesResponse, error) {
	return s.rzdClient.GetTrainCarriages(ctx, params)
}

// GetTrainStationList реализует метод получения списка станций поезда
func (s *mainService) GetTrainStationList(ctx context.Context, params domain.GetTrainStationListParams) (domain.TrainStationListResponse, error) {
	return s.rzdClient.GetTrainStationList(ctx, params)
}

// GetStationCode реализует метод получения кодов станций
func (s *mainService) GetStationCode(ctx context.Context, params domain.GetStationCodeParams) ([]domain.StationCode, error) {
	return s.rzdClient.GetStationCode(ctx, params)
}
