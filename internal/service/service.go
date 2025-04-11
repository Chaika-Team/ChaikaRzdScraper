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

// GetTrainRoutes получение маршрутов поездов
func (s *mainService) GetTrainRoutes(ctx context.Context, params domain.GetTrainRoutesParams) ([]domain.TrainRoute, error) {
	return s.rzdClient.GetTrainRoutes(ctx, params)
}

// GetTrainCarriages получение информации о вагонах
func (s *mainService) GetTrainCarriages(ctx context.Context, params domain.GetTrainCarriagesParams) ([]domain.Car, error) {
	return s.rzdClient.GetTrainCarriages(ctx, params)
}

// SearchStation получение кодов станций по поисковому запросу
func (s *mainService) SearchStation(ctx context.Context, params domain.SearchStationParams) ([]domain.Station, error) {
	return s.rzdClient.SearchStation(ctx, params)
}
