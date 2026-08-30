package testutils

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

type PostgresTestContainer struct {
	Container        *postgres.PostgresContainer
	DB               *sql.DB
	ConnectionString string
}

func StartPostgres(ctx context.Context) (*PostgresTestContainer, error) {

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil, errors.New("failed to resolve test helper path")
	}

	projectRoot := filepath.Clean(
		filepath.Join(filepath.Dir(filename), "../../.."),
	)

	migrationFiles, err := filepath.Glob(
		filepath.Join(projectRoot, "migration", "*.up.sql"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to find migrations: %w", err)
	}

	if len(migrationFiles) == 0 {
		return nil, errors.New("migration files not found")
	}

	pgContainer, err := postgres.Run(
		ctx,
		"postgres:15-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		postgres.WithOrderedInitScripts(migrationFiles...),
		postgres.BasicWaitStrategies(),
	)

	connectionString, err := pgContainer.ConnectionString(
		ctx,
		"sslmode=disable",
	)
	if err != nil {
		_ = testcontainers.TerminateContainer(pgContainer)

		return nil, fmt.Errorf("failed to get connection string: %w", err)
	}

	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		_ = testcontainers.TerminateContainer(pgContainer)

		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		_ = testcontainers.TerminateContainer(pgContainer)

		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresTestContainer{
		Container:        pgContainer,
		DB:               db,
		ConnectionString: connectionString,
	}, nil
}

func (c *PostgresTestContainer) Close() error {
	var errs []error

	if err := c.DB.Close(); err != nil {
		errs = append(errs, fmt.Errorf("failed to close databse: %w", err))
	}

	if err := testcontainers.TerminateContainer(c.Container); err != nil {
		errs = append(errs, fmt.Errorf("failed to terminate postgres: %w", err))
	}

	return errors.Join(errs...)
}
