package tests

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	mpostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/subvisual/fidl/bank"
	"github.com/subvisual/fidl/bank/postgres"
	"github.com/subvisual/fidl/tests/setup"
)

// nolint:gochecknoglobals
var (
	db                *postgres.DB
	migr              *migrate.Migrate
	bankFqdn          string
	bankPort          int
	proxyPrice        = "1 FIL"
	bankWalletAddress string
	containerDeadline uint = 3000
)

func TestMain(m *testing.M) {
	cfgFilePath := "../etc/bank.ini.example"
	cfg := bank.LoadConfiguration(cfgFilePath)

	bankFqdn = cfg.HTTP.Fqdn
	bankPort = cfg.HTTP.Port
	bankWalletAddress = cfg.Wallet.Address.String()

	db = &postgres.DB{}
	pool, resource, err := setup.PostgresContainer(db, &cfg, containerDeadline)
	if err != nil {
		log.Fatalf("could not setup docker: %v", err)
	}

	go func() {
		cfg.Logger.Level = "FATAL"
		setup.Server(cfg, db)
	}()

	res, err := ServerHealthcheck(cfg)
	if res.Body.Close(); err != nil {
		log.Fatalf("failed to check server health: %v", err)
	}

	driver, err := mpostgres.WithInstance(db.DB.DB, &mpostgres.Config{})
	if err != nil {
		log.Fatalf("could not create migrate db instance: %v", err)
	}

	migr, err = migrate.NewWithDatabaseInstance(
		"file://../bank/postgres/migrations",
		"postgres", driver)
	if err != nil {
		log.Fatalf("could not get migrations: %v", err)
	}

	defer func() {
		if err := pool.Purge(resource); err != nil {
			log.Fatalf("could not purge resource: %s", err)
		}
	}()

	m.Run()
}

func ServerHealthcheck(cfg bank.Config) (*http.Response, error) {
	endpoint, err := url.JoinPath("http://"+cfg.HTTP.Fqdn+":"+fmt.Sprint(cfg.HTTP.Port), "/api/v1/healthcheck")
	if err != nil {
		return nil, fmt.Errorf("failed joining endpoint path: %w", err)
	}

	timeout := 5
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}

	retries := 5
	client := http.DefaultClient
	res, err := client.Do(req)
	for err != nil {
		if retries > 1 {
			retries--
			time.Sleep(1 * time.Second)
			res, err = http.DefaultClient.Do(req)

			continue
		}

		return nil, fmt.Errorf("failed to boot server %d times: %w", retries, err)
	}

	return res, nil
}
