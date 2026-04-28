package main

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jesse-rb/pmt-excel-fn-go/protogen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type httpServer struct {
	addr string
}

func NewHTTPServer(addr string) *httpServer {
	return &httpServer{addr: addr}
}

func (s *httpServer) Run(ctx context.Context) error {
	router := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := protogen.RegisterPMTServiceHandlerFromEndpoint(ctx, router, "localhost:9090", opts)
	if err != nil {
		return err
	}

	slog.Info("Starting pmt http server and proxy calls to grpc server", "on", s.addr)

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	return http.ListenAndServe(s.addr, router)
}
