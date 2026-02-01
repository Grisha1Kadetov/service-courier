package testutils

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupDB(t *testing.T) (*pgxpool.Pool, func()) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:15",
		ExposedPorts: []string{"0:5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "user",
			"POSTGRES_PASSWORD": "pass",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForAll(
			wait.ForListeningPort("5432/tcp"),
			wait.ForLog("database system is ready to accept connections"),
		),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	host, err := container.Host(ctx)
	require.NoError(t, err)

	port, err := container.MappedPort(ctx, "5432")
	require.NoError(t, err)

	dsn := fmt.Sprintf(
		"postgres://user:pass@%s:%s/testdb",
		host, port.Port(),
	)

	var pool *pgxpool.Pool
	for i := 0; i < 10; i++ {
		pool, err = pgxpool.New(ctx, dsn)
		if err == nil {
			err = pool.Ping(ctx)
			if err == nil {
				break
			}
		}
		time.Sleep(time.Second)
	}
	require.NoError(t, err)

	db, err := goose.OpenDBWithDriver("pgx", dsn)
	require.NoError(t, err)
	d, err := os.Getwd()
	require.NoError(t, err)
	err = goose.Up(db, FindModuleRoot(d)+"/migration")
	require.NoError(t, err)

	cleanup := func() {
		pool.Close()
		_ = container.Terminate(ctx)
	}

	return pool, cleanup
}
