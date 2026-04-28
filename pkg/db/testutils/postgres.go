package testutils

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jesse-rb/pmt-excel-fn-go/pkg/db"
	"github.com/ory/dockertest/v4"
)

func SetupTestPostgres(t *testing.T, ctx context.Context) (*db.Postgres, error) {
	pool := dockertest.NewPoolT(t, "")

	postgresResource := pool.RunT(t, "postgres",
		dockertest.WithTag("18"),
		dockertest.WithEnv([]string{
			"POSTGRES_PASSWORD=secret",
			"POSTGRES_DB=testdb",
		}),
	)

	hostPort := postgresResource.GetHostPort("5432/tcp")
	dsn := fmt.Sprintf("postgres://postgres:secret@%s/testdb", hostPort)

	db, err := db.NewPostgres(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to init test postgres: %w", err)
	}

	// Wait for PostgreSQL to be ready
	err = pool.Retry(ctx, 30*time.Second, func() error {
		return db.Pool.Ping(ctx)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to test postgres: %w", err)
	}

	db.RunMigrations()

	return db, err
}
