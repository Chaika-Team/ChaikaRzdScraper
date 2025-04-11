package grpc

import (
	"context"
	"github.com/go-kit/kit/endpoint"

	"github.com/Chaika-Team/rzd-api/internal/transports/grpc/mappers"

	"github.com/Chaika-Team/rzd-api/internal/domain"
	"github.com/Chaika-Team/rzd-api/internal/service"
	"github.com/Chaika-Team/rzd-api/internal/transports/grpc/pb"
)

// Endpoints собраны для gRPC сервиса.
type Endpoints struct {
	GetTrainRoutes    endpoint.Endpoint
	GetTrainCarriages endpoint.Endpoint
	SearchStation     endpoint.Endpoint
}

// MakeEndpoints создаёт эндпоинты из сервиса.
func MakeEndpoints(svc service.Service) Endpoints {
	return Endpoints{
		GetTrainRoutes:    makeGetTrainRoutesEndpoint(svc),
		GetTrainCarriages: makeGetTrainCarriagesEndpoint(svc),
		SearchStation:     makeSearchStationEndpoint(svc),
	}
}

func makeGetTrainRoutesEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*pb.GetTrainRoutesRequest)
		params := domain.GetTrainRoutesParams{
			FromCode:   int(req.FromCode),
			ToCode:     int(req.ToCode),
			Direction:  domain.Direction(req.Direction),
			TrainType:  domain.TrainSearchType(req.TrainType),
			CheckSeats: req.CheckSeats,
			FromDate:   mappers.ParseDateRequest(req.FromDate),
			WithChange: req.WithChange,
		}
		routes, err := svc.GetTrainRoutes(ctx, params)
		if err != nil {
			return nil, err
		}
		return mappers.MapTrainRoutesToPb(routes), nil
	}
}

func makeGetTrainCarriagesEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*pb.GetTrainCarriagesRequest)
		params := domain.GetTrainCarriagesParams{
			TrainNumber: req.TrainNumber,
			Direction:   domain.Direction(req.Direction),
			FromCode:    int(req.FromCode),
			FromTime:    mappers.ParseTimeRequest(req.FromTime),
			ToCode:      int(req.ToCode),
		}
		cars, err := svc.GetTrainCarriages(ctx, params)
		if err != nil {
			return nil, err
		}
		return mappers.MapTrainCarriagesToPb(cars), nil
	}
}

func makeSearchStationEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*pb.SearchStationRequest)
		params := domain.SearchStationParams{
			Query:       req.Query,
			CompactMode: req.CompactMode,
		}
		stations, err := svc.SearchStation(ctx, params)
		if err != nil {
			return nil, err
		}
		return mappers.MapStationsToPb(stations), nil
	}
}
