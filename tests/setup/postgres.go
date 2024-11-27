package setup

import (
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file" // needed for file driver for migrations
	"github.com/jmoiron/sqlx"
	dockertest "github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/subvisual/fidl/bank/postgres"
)

func PostgresContainer(db *postgres.DB) (*dockertest.Pool, *dockertest.Resource, string, error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, nil, "", fmt.Errorf("could not connect to docker: %w", err)
	}

	user := "user"
	repo := "postgres"
	version := "17"
	pwd := "secret"
	dbname := "fidl-postgres-dev"

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Name:       "fidl-postgres-testing",
		Repository: repo,
		Tag:        version,
		Env: []string{
			fmt.Sprintf("POSTGRES_PASSWORD=%s", pwd),
			fmt.Sprintf("POSTGRES_USER=%s", user),
			fmt.Sprintf("POSTGRES_DB=%s", dbname),
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})

	if err != nil {
		return nil, nil, "", fmt.Errorf("could not start resource: %w", err)
	}

	if err := resource.Expire(30); err != nil {
		return nil, nil, "", fmt.Errorf("could not set expiration to resource: %w", err)
	}

	var dbx *sqlx.DB
	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseURL := fmt.Sprintf("%s://%s:%s@%s/%s?sslmode=disable", repo, user, pwd, hostAndPort, dbname)

	pool.MaxWait = 30 * time.Second
	if err = pool.Retry(func() error {
		dbx, err = sqlx.Open(repo, databaseURL)
		if err != nil {
			return fmt.Errorf("could not connect to database: %w", err)
		}

		return dbx.Ping()
	}); err != nil {
		return nil, nil, "", fmt.Errorf("could not connect to docker: %w", err)
	}

	db.DB = dbx

	return pool, resource, "", nil
}

func RunMigrations(kind string, migr *migrate.Migrate) error {
	switch kind {
	case "UP":
		if err := migr.Up(); err != nil {
			return fmt.Errorf("failed to run up migrations: %w", err)
		}
	case "DOWN":
		if err := migr.Down(); err != nil {
			return fmt.Errorf("failed to run down migrations: %w", err)
		}
	default:
		return fmt.Errorf("wrong given kind of migration: %s", kind)
	}

	return nil
}
