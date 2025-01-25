// internal/infrastructure/grpc/server.go
package grpc

import (
	"context"

	"github.com/Chaika-Team/rzd-api/internal/domain"
	pb "github.com/Chaika-Team/rzd-api/internal/infrastructure/grpc/pb"
	"github.com/Chaika-Team/rzd-api/internal/usecase"
)

type RZDServiceServer struct {
	pb.UnimplementedRZDServiceServer
	RouteUseCase    *usecase.RouteUseCase
	CarriageUseCase *usecase.CarriageUseCase
	StationUseCase  *usecase.StationUseCase
}

func NewRZDServiceServer(routeUC *usecase.RouteUseCase, carriageUC *usecase.CarriageUseCase, stationUC *usecase.StationUseCase) *RZDServiceServer {
	return &RZDServiceServer{
		RouteUseCase:    routeUC,
		CarriageUseCase: carriageUC,
		StationUseCase:  stationUC,
	}
}

func (s *RZDServiceServer) GetTrainRoutes(ctx context.Context, req *pb.GetTrainRoutesRequest) (*pb.GetTrainRoutesResponse, error) {
	params := domain.GetTrainRoutesParams{
		Code0:      req.GetCode0(),
		Code1:      req.GetCode1(),
		Dir:        req.GetDir(),
		Tfl:        req.GetTfl(),
		CheckSeats: req.GetCheckSeats(),
		Dt0:        req.GetDt0(),
		Md:         req.GetMd(),
	}

	routes, err := s.RouteUseCase.GetTrainRoutes(ctx, params)
	if err != nil {
		return nil, err
	}

	var pbRoutes []*pb.TrainRoute
	for _, route := range routes {
		pbRoutes = append(pbRoutes, &pb.TrainRoute{
			Route0:    route.Route0,
			Route1:    route.Route1,
			Date0:     route.DepartureDate,
			Time0:     route.DepartureTime,
			Number:    route.Number,
			From:      route.From,
			Where:     route.Where,
			Date:      route.ArrivalDate,
			FromCode:  route.FromCode,
			WhereCode: route.WhereCode,
			Time1:     route.ArrivalTime,
			TimeInWay: route.TimeInWay,
			Brand:     route.Brand,
			Carrier:   route.Carrier,
		})
	}

	return &pb.GetTrainRoutesResponse{
		Routes: pbRoutes,
	}, nil
}

func (s *RZDServiceServer) GetTrainRoutesReturn(ctx context.Context, req *pb.GetTrainRoutesReturnRequest) (*pb.GetTrainRoutesReturnResponse, error) {
	params := domain.GetTrainRoutesReturnParams{
		Code0:      req.GetCode0(),
		Code1:      req.GetCode1(),
		Dir:        req.GetDir(),
		Tfl:        req.GetTfl(),
		CheckSeats: req.GetCheckSeats(),
		Dt0:        req.GetDt0(),
		Dt1:        req.GetDt1(),
	}

	forward, back, err := s.RouteUseCase.GetTrainRoutesReturn(ctx, params)
	if err != nil {
		return nil, err
	}

	var pbForward []*pb.TrainRoute
	for _, route := range forward {
		pbForward = append(pbForward, &pb.TrainRoute{
			Route0:    route.Route0,
			Route1:    route.Route1,
			Date0:     route.DepartureDate,
			Time0:     route.DepartureTime,
			Number:    route.Number,
			From:      route.From,
			Where:     route.Where,
			Date:      route.ArrivalDate,
			FromCode:  route.FromCode,
			WhereCode: route.WhereCode,
			Time1:     route.ArrivalTime,
			TimeInWay: route.TimeInWay,
			Brand:     route.Brand,
			Carrier:   route.Carrier,
		})
	}

	var pbBack []*pb.TrainRoute
	for _, route := range back {
		pbBack = append(pbBack, &pb.TrainRoute{
			Route0:    route.Route0,
			Route1:    route.Route1,
			Date0:     route.DepartureDate,
			Time0:     route.DepartureTime,
			Number:    route.Number,
			From:      route.From,
			Where:     route.Where,
			Date:      route.ArrivalDate,
			FromCode:  route.FromCode,
			WhereCode: route.WhereCode,
			Time1:     route.ArrivalTime,
			TimeInWay: route.TimeInWay,
			Brand:     route.Brand,
			Carrier:   route.Carrier,
		})
	}

	return &pb.GetTrainRoutesReturnResponse{
		Forward: pbForward,
		Back:    pbBack,
	}, nil
}

func (s *RZDServiceServer) GetTrainCarriages(ctx context.Context, req *pb.GetTrainCarriagesRequest) (*pb.GetTrainCarriagesResponse, error) {
	params := domain.GetTrainCarriagesParams{
		Code0: req.GetCode0(),
		Code1: req.GetCode1(),
		Tnum0: req.GetTnum0(),
		Time0: req.GetTime0(),
		Dt0:   req.GetDt0(),
		Dir:   req.GetDir(),
	}

	carriages, err := s.CarriageUseCase.GetTrainCarriages(ctx, params)
	if err != nil {
		return nil, err
	}

	var pbCars []*pb.Carriage
	for _, car := range carriages.Cars {
		var pbSeats []*pb.Seat
		for _, seat := range car.Seats {
			pbSeats = append(pbSeats, &pb.Seat{
				Places: seat.Places,
				Tariff: seat.Tariff,
				Type:   seat.Type,
				Free:   seat.Free,
				Label:  seat.Label,
			})
		}

		pbCars = append(pbCars, &pb.Carriage{
			Cnumber: car.CNumber,
			Type:    car.Type,
			TypeLoc: car.TypeLoc,
			ClsType: car.ClsType,
			Tariff:  car.Tariff,
			Seats:   pbSeats,
		})
	}

	return &pb.GetTrainCarriagesResponse{
		Response: &pb.TrainCarriagesResponse{
			Cars:           pbCars,
			FunctionBlocks: carriages.FunctionBlocks,
			Schemes:        carriages.Schemes,
			Companies:      carriages.Companies,
		},
	}, nil
}

func (s *RZDServiceServer) GetTrainStationList(ctx context.Context, req *pb.GetTrainStationListRequest) (*pb.GetTrainStationListResponse, error) {
	params := domain.GetTrainStationListParams{
		TrainNumber: req.GetTrainNumber(),
		DepDate:     req.GetDepDate(),
	}

	stations, err := s.StationUseCase.GetTrainStationList(ctx, params)
	if err != nil {
		return nil, err
	}

	var pbRoutes []*pb.Route
	for _, route := range stations.Routes {
		pbRoutes = append(pbRoutes, &pb.Route{
			Station:       route.Station,
			ArrivalTime:   route.ArrivalTime,
			DepartureTime: route.DepartureTime,
			// Добавьте другие поля по необходимости
		})
	}

	return &pb.GetTrainStationListResponse{
		Response: &pb.TrainStationListResponse{
			Train: &pb.TrainInfo{
				Number: stations.Train.Number,
				// Добавьте другие поля
			},
			Routes: pbRoutes,
		},
	}, nil
}

func (s *RZDServiceServer) GetStationCode(ctx context.Context, req *pb.GetStationCodeRequest) (*pb.GetStationCodeResponse, error) {
	params := domain.GetStationCodeParams{
		StationNamePart: req.GetStationNamePart(),
		CompactMode:     req.GetCompactMode(),
	}

	stationCodes, err := s.StationUseCase.GetStationCode(ctx, params)
	if err != nil {
		return nil, err
	}

	var pbStations []*pb.StationCode
	for _, sc := range stationCodes {
		pbStations = append(pbStations, &pb.StationCode{
			Station: sc.Station,
			Code:    sc.Code,
		})
	}

	return &pb.GetStationCodeResponse{
		Stations: pbStations,
	}, nil
}
