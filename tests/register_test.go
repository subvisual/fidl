package tests

import (
	"log"
	"testing"

	"github.com/subvisual/fidl/proxy"
	"github.com/subvisual/fidl/tests/setup"
)

func TestRegister(t *testing.T) { // nolint:paralleltest
	if err := setup.RunMigrations("UP", migr); err != nil {
		log.Fatalf("could not run up migrations: %v", err)
	}

	cfg, _ := setup.Proxy("1 FIL", localhost, bankPort, upstreamPort)

	if err := proxy.Register(cfg); err != nil {
		t.Log("failed to register proxy", err)
		t.Fail()
	}

	if err := setup.RunMigrations("DOWN", migr); err != nil {
		log.Fatalf("could not run down migrations: %v", err)
	}
}