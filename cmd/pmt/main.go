package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/jesse-rb/pmt-excel-fn-go/pkg/db"
	"github.com/jesse-rb/pmt-excel-fn-go/pkg/pmt"
)

func main() {
	ctx := context.Background()

	// Connect to DB
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&TimeZone=UTC", dbUser, dbPass, dbHost, dbPort, dbName)

	db, err := db.NewPostgres(context.Background(), dsn)
	defer db.Close()
	if err != nil {
		slog.Error("error creating new postgres connection", "err", err)
	}

	// Check DB connection before continuing
	for i := 0; i < 10; i++ {
		if db.Pool.Ping(ctx) == nil {
			break
		}
		time.Sleep(time.Second)
	}

	// For local dev (for convenience), we can run DB migrations on startup
	db.RunMigrations()

	pg := pmt.NewPostgresRepository(db.Pool)

	// // Setup HTTP handler and routes
	// httpHandler := pmt.NewHTTPHandler(pg)
	// mux := http.NewServeMux()
	// mux.HandleFunc("/pmt", httpHandler.HandlePMT)
	//
	// // Listen and serve HTTP
	// httpServer := &http.Server{
	// 	Addr:    ":8080",
	// 	Handler: mux,
	// }
	// defer httpServer.Close()
	// httpServer.ListenAndServe()

	grpcServer := NewGRPCServer(":9090", pg)
	go (func() {
		err = grpcServer.Run() // Blocking, so we run in a go-rotine
		if err != nil {
			slog.Error("error starting pmt grpc server", "err", err)
		}
	})()

	httpProxy := NewHTTPServer(":8080")
	err = httpProxy.Run(ctx)
	if err != nil {
		slog.Error("error starting pmt http grpc proxy server", "err", err)
	}
}
