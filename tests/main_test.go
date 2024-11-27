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
	fhttp "github.com/subvisual/fidl/http"
	"github.com/subvisual/fidl/tests/setup"
	"github.com/subvisual/fidl/types"
)

// nolint:gochecknoglobals
var (
	db           *postgres.DB
	migr         *migrate.Migrate
	localhost    = "localhost"
	bankPort     = 8090
	upstreamPort = 7777
	proxyPrice   = "1 FIL"
)

func TestMain(m *testing.M) {
	db = &postgres.DB{}
	pool, resource, dbDSN, err := setup.PostgresContainer(db)
	if err != nil {
		log.Fatalf("could not setup docker: %v", err)
	}

	bankAddress, _ := types.NewAddressFromString("f1qbvbikeuozxgoop5bc7nkcokapslxxdy2gucfqa")
	escrowAddress, _ := types.NewAddressFromString("f1bpzzps6xxcxubq6idqravdcisdpzqemnahxdeoq")

	cfg := bank.Config{
		Env: "testing",
		Logger: fhttp.Logger{
			Path:  "logs/bank.log",
			Level: "FATAL",
		},
		Db: bank.Db{
			Dsn:          dbDSN,
			MaxOpenConns: 25,
			MaxIdleConns: 25,
			MaxIdleTime:  "15m",
		},
		HTTP: fhttp.HTTP{
			Addr:            "127.0.0.1",
			Fqdn:            localhost,
			Port:            bankPort,
			ListenPort:      bankPort,
			ReadTimeout:     15,
			WriteTimeout:    15,
			ShutdownTimeout: 10,
			TLS:             false,
		},
		Wallet: bank.Wallet{
			Address: bankAddress,
		},
		Escrow: bank.Escrow{
			Address:  escrowAddress,
			Deadline: "24h",
		},
	}

	go func() {
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
