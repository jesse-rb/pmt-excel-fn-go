package main

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/jesse-rb/pmt-excel-fn-go/pkg/pmt"
	"google.golang.org/grpc"
)

type grpcServer struct {
	addr  string
	store pmt.Store
}

func NewGRPCServer(addr string, store pmt.Store) *grpcServer {
	return &grpcServer{addr: addr, store: store}
}

func (s *grpcServer) Run() error {
	// Create a TCP listener to pass to grpc server
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("failed to listen for grpc server: %w", err)
	}

	grpcServer := grpc.NewServer()

	// Register our gRPC service
	pmtService := pmt.NewService()
	pmt.NewGRPCHandler(grpcServer, pmtService, s.store)

	slog.Info("Starting pmt grpc server", "on", s.addr)

	return grpcServer.Serve(listener)
}
