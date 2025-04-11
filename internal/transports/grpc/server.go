package grpc

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/Chaika-Team/rzd-api/internal/transports/grpc/pb"

	"google.golang.org/grpc"
)

// Server реализует pb.RzdServiceServer.
type Server struct {
	endpoints Endpoints
	pb.UnimplementedRzdServiceServer
}

// NewGRPCServer создаёт новый gRPC сервер.
func NewGRPCServer(endpoints Endpoints) *Server {
	return &Server{endpoints: endpoints}
}

func (s *Server) GetTrainRoutes(ctx context.Context, req *pb.GetTrainRoutesRequest) (*pb.GetTrainRoutesResponse, error) {
	response, err := s.endpoints.GetTrainRoutes(ctx, req)
	if err != nil {
		return nil, err
	}
	resp, ok := response.(*pb.GetTrainRoutesResponse)
	if !ok {
		return nil, fmt.Errorf("unexpected response type: %T", response)
	}
	return resp, nil
}

func (s *Server) GetTrainCarriages(ctx context.Context, req *pb.GetTrainCarriagesRequest) (*pb.GetTrainCarriagesResponse, error) {
	response, err := s.endpoints.GetTrainCarriages(ctx, req)
	if err != nil {
		return nil, err
	}
	resp, ok := response.(*pb.GetTrainCarriagesResponse)
	if !ok {
		return nil, fmt.Errorf("unexpected response type: %T", response)
	}
	return resp, nil
}

func (s *Server) SearchStation(ctx context.Context, req *pb.SearchStationRequest) (*pb.SearchStationResponse, error) {
	response, err := s.endpoints.SearchStation(ctx, req)
	if err != nil {
		return nil, err
	}
	resp, ok := response.(*pb.SearchStationResponse)
	if !ok {
		return nil, fmt.Errorf("unexpected response type: %T", response)
	}
	return resp, nil
}

// StartGRPCServer запускает gRPC-сервер и возвращает grpc.Server для управления его остановкой.
func StartGRPCServer(addr string, srv *Server) (*grpc.Server, net.Listener, error) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, nil, err
	}
	grpcServer := grpc.NewServer()
	pb.RegisterRzdServiceServer(grpcServer, srv)
	log.Printf("gRPC server listening on %s", addr)
	return grpcServer, listener, nil
}
