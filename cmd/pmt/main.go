package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/jesse-rb/pmt-excel-fn-go/pkg/db"
	"github.com/jesse-rb/pmt-excel-fn-go/pkg/pmt"
)

func main() {
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
		// db, err := sql.Open("postgres", dsn)
		if db.Pool.Ping(context.Background()) == nil {
			break
		}
		time.Sleep(time.Second)
	}

	// For local dev (for convenience), we can run DB migrations on startup
	db.RunMigrations()

	// Listen and serve HTTP server
	httpHandler := pmt.NewHTTPHandler(nil)
	mux := http.NewServeMux()
	mux.HandleFunc("/pmt", httpHandler.HandlePMT)
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	httpServer.ListenAndServe()
}
