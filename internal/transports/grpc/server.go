package grpc

import (
	"context"
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
	return response.(*pb.GetTrainRoutesResponse), nil
}

func (s *Server) GetTrainCarriages(ctx context.Context, req *pb.GetTrainCarriagesRequest) (*pb.GetTrainCarriagesResponse, error) {
	response, err := s.endpoints.GetTrainCarriages(ctx, req)
	if err != nil {
		return nil, err
	}
	return response.(*pb.GetTrainCarriagesResponse), nil
}

func (s *Server) SearchStation(ctx context.Context, req *pb.SearchStationRequest) (*pb.SearchStationResponse, error) {
	response, err := s.endpoints.SearchStation(ctx, req)
	if err != nil {
		return nil, err
	}
	return response.(*pb.SearchStationResponse), nil
}

// StartGRPCServer Запуск gRPC-сервера.
func StartGRPCServer(addr string, srv *Server) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	pb.RegisterRzdServiceServer(grpcServer, srv)
	log.Printf("gRPC server listening on %s", addr)
	return grpcServer.Serve(listener)
}
