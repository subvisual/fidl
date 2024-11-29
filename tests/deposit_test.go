package tests

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/subvisual/fidl/cli"
	"github.com/subvisual/fidl/tests/setup"
)

func TestDeposit(t *testing.T) { // nolint:paralleltest
	if err := setup.RunMigrations("UP", migr); err != nil {
		t.Fatalf("could not run up migrations: %v", err)
	}

	cfg, cl, ki, err := setup.CLI()
	if err != nil {
		t.Fatalf("could not setup CLI info: %v", err)
	}

	bankEndpoint := url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s:%d", bankFqdn, bankPort),
	}

	var tests = []struct {
		address  string
		amount   string
		expected string
	}{
		{bankEndpoint.String(), proxyPrice, "1 FIL"},
		{bankEndpoint.String(), proxyPrice, "2 FIL"},
		{bankEndpoint.String(), proxyPrice, "3 FIL"},
		{bankEndpoint.String(), proxyPrice, "4 FIL"},
		{bankEndpoint.String(), proxyPrice, "5 FIL"},
	}

	for _, test := range tests {
		opts := cli.DepositOptions{
			Amount:      test.amount,
			BankAddress: test.address,
		}

		if err := cl.Validate.Struct(opts); err != nil {
			t.Errorf("failed to validate: %v", err)
		}

		res, err := cli.Deposit(ki, cfg.Wallet.Address, cfg.Route.Deposit, opts)
		if err != nil {
			t.Errorf("failed to deposit: %v", err)
		}

		assert.Equal(t, res.Status, "success")
		assert.Equal(t, res.Data.FIL.String(), test.expected)
	}

	if err := setup.RunMigrations("DOWN", migr); err != nil {
		t.Fatalf("could not run down migrations: %v", err)
	}
}
