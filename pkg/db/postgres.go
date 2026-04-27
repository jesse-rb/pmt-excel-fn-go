package db

import (
	"context"
	"embed"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrations embed.FS

type Postgres struct {
	Pool *pgxpool.Pool
}

func NewPostgres(ctx context.Context, dsn string) (*Postgres, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create new pgxpool: %w", err)
	}
	return &Postgres{Pool: pool}, nil
}

func (db *Postgres) Close() {
	if db.Pool == nil {
		slog.Warn("postgres close conn: db pool is nil")
	} else {
		db.Pool.Close()
	}
}

func (db *Postgres) RunMigrations() {
	if err := goose.SetDialect("postgres"); err != nil {
		slog.Error("error setting goose dialect", "err", err)
	}
	goose.SetBaseFS(migrations)

	sqldb := stdlib.OpenDBFromPool(db.Pool)
	defer (func() {
		if err := sqldb.Close(); err != nil {
			slog.Error("error getting *sql.DB from *pgxpool.Pool", "err", err)
		}
	})()
	if err := goose.Up(sqldb, "migrations"); err != nil {
		slog.Error("error running migrations", "err", err)
	}

}
