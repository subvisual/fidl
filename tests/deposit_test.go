package tests

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/subvisual/fidl/cli"
	"github.com/subvisual/fidl/tests/setup"
)

func TestDeposit(t *testing.T) { // nolint:paralleltest
	if err := setup.RunMigrations("UP", migr); err != nil {
		t.Fatalf("could not run up migrations: %v", err)
	}

	cfg, _, ki, err := setup.CLI()
	if err != nil {
		t.Fatalf("could not setup CLI info: %v", err)
	}

	bankAddress := fmt.Sprintf("http://%s:%d", localhost, bankPort)

	var tests = []struct {
		address  string
		amount   string
		expected string
	}{
		{bankAddress, "1 FIL", "1 FIL"},
		{bankAddress, "1 FIL", "2 FIL"},
		{bankAddress, "1 FIL", "3 FIL"},
		{bankAddress, "1 FIL", "4 FIL"},
		{bankAddress, "1 FIL", "5 FIL"},
		{bankAddress, "1 FIL", "6 FIL"},
	}

	for _, test := range tests {
		opts := cli.DepositOptions{
			Amount:      test.amount,
			BankAddress: test.address,
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
