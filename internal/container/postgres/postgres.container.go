package postgres

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/NekKkMirror/medods-tz.git/config"
	"github.com/cenkalti/backoff/v4"
	"github.com/docker/go-connections/nat"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Options holds configuration for the PostgreSQL container.
type Options struct {
	Database  string
	Host      string
	Port      nat.Port
	UserName  string
	Password  string
	ImageName string
	Name      string
	Tag       string
	Timeout   time.Duration
}

// init sets up environment variables for PostgreSQL connection.
func init() {
	os.Setenv("POSTGRES_USER", "test")
	os.Setenv("POSTGRES_PASSWORD", "test")
	os.Setenv("POSTGRES_HOST", "localhost")
	os.Setenv("POSTGRES_DB", "test")
}

// Start initializes a PostgreSQL container and returns a sqlx.DB instance and any error occurred.
func Start(ctx context.Context) (*sqlx.DB, func(), error) {
	options := getDefaultPostgresOptions()
	containerReq := getContainerRequest(options)

	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: containerReq,
		Started:          true,
	})
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to start PostgreSQL container")
	}

	mappedPort, err := postgresContainer.MappedPort(ctx, options.Port)
	if err != nil {
		_ = postgresContainer.Terminate(ctx)
		return nil, nil, errors.Wrap(err, "failed to get exposed container port")
	}
	os.Setenv("POSTGRES_PORT", mappedPort.Port())

	DB, err := createDBConnection(ctx, postgresContainer, options)
	if err != nil {
		_ = postgresContainer.Terminate(ctx)
		return nil, nil, errors.Wrap(err, "failed to create DB connection")
	}

	if err := createTables(DB); err != nil {
		_ = postgresContainer.Terminate(ctx)
		return nil, nil, errors.Wrap(err, "failed to create tables")
	}

	return DB, func() {
		if err := postgresContainer.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %s", err)
		}
	}, nil
}

// getDefaultPostgresOptions returns the default configuration for PostgreSQL container.
func getDefaultPostgresOptions() *Options {
	port, err := nat.NewPort("tcp", "5432")
	if err != nil {
		panic(errors.Wrap(err, "failed to create new port"))
	}

	return &Options{
		Database:  "test",
		Port:      port,
		Host:      "localhost",
		UserName:  "test",
		Password:  "test",
		Tag:       "latest",
		ImageName: "postgres",
		Name:      "postgresql-testcontainer",
		Timeout:   5 * time.Minute,
	}
}

// getContainerRequest builds and returns a testcontainers.ContainerRequest using the provided options.
func getContainerRequest(opts *Options) testcontainers.ContainerRequest {
	return testcontainers.ContainerRequest{
		Image:        fmt.Sprintf("%s:%s", opts.ImageName, opts.Tag),
		ExposedPorts: []string{opts.Port.Port()},
		WaitingFor:   wait.ForListeningPort(opts.Port),
		Env: map[string]string{
			"POSTGRES_DB":       opts.Database,
			"POSTGRES_PASSWORD": opts.Password,
			"POSTGRES_USER":     opts.UserName,
		},
	}
}

// createDBConnection establishes a sqlx connection using provided PostgreSQL container and options.
func createDBConnection(ctx context.Context, container testcontainers.Container, opts *Options) (*sqlx.DB, error) {
	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = 10 * time.Second
	const maxRetries = 5

	var (
		DB  *sqlx.DB
		err error
	)

	err = backoff.Retry(func() error {
		opts.Host, err = container.Host(ctx)
		if err != nil {
			return errors.Wrap(err, "failed to get container host")
		}

		opts.Port, err = container.MappedPort(ctx, opts.Port)
		if err != nil {
			return errors.Wrap(err, "failed to get exposed container port")
		}

		DB, err = config.ConnectDB()
		return err
	}, backoff.WithMaxRetries(bo, maxRetries))

	if err != nil {
		return nil, errors.Wrap(err, "failed to create connection after retries")
	}

	return DB, nil
}

// createTables creates the necessary tables in the provided PostgreSQL database.
func createTables(db *sqlx.DB) error {
	createUserTable := `create table users (
    id uuid not null primary key,
    email varchar(255) not null unique
	);`
	_, err := db.Exec(createUserTable)
	if err != nil {
		return errors.Wrap(err, "failed to create users table")
	}
	_, err = db.Exec("INSERT INTO users (id, email) VALUES ('3100b0c6-c6cc-4edf-a9a8-444990d0547d', 'user1@example.com');")
	if err != nil {
		return errors.Wrap(err, "failed to insert user1")
	}

	createRefreshTokenTable := `create table refresh_tokens(
    user_id uuid not null primary key constraint fk_user references users on delete cascade,
    token_hash varchar(255) not null,
    client_ip  varchar(255) not null
	);`
	_, err = db.Exec(createRefreshTokenTable)
	if err != nil {
		return errors.Wrap(err, "failed to create refresh_tokens table")
	}

	return nil
}
