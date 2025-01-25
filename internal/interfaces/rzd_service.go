// internal/interfaces/rzd_service.go
package interfaces

import (
	"github.com/Chaika-Team/rzd-api/internal/infrastructure/grpc"
	"github.com/Chaika-Team/rzd-api/internal/usecase"
)

func NewRZDServiceServer(routeUC *usecase.RouteUseCase, carriageUC *usecase.CarriageUseCase, stationUC *usecase.StationUseCase) *grpc.RZDServiceServer {
	return grpc.NewRZDServiceServer(routeUC, carriageUC, stationUC)
}
