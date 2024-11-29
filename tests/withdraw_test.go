package tests

import (
	"fmt"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/subvisual/fidl/cli"
	"github.com/subvisual/fidl/tests/setup"
)

func TestWithdraw(t *testing.T) { // nolint:paralleltest
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

	depositOpts := cli.DepositOptions{
		Amount:      "5 FIL",
		BankAddress: bankEndpoint.String(),
	}

	if err := cl.Validate.Struct(depositOpts); err != nil {
		t.Errorf("failed to validate: %v", err)
	}

	res, err := cli.Deposit(ki, cfg.Wallet.Address, cfg.Route.Deposit, depositOpts)
	if err != nil {
		t.Errorf("failed to deposit: %v", err)
	}

	assert.Equal(t, res.Status, "success")
	assert.Equal(t, res.Data.FIL.String(), "5 FIL")

	destinationAddress := "f135aafnv6wnlderpanbbfwwc3zxnxzomsphqnnfq"

	var tests = []struct {
		address     string
		destination string
		amount      string
		expected    string
	}{
		{bankEndpoint.String(), destinationAddress, proxyPrice, "4 FIL"},
		{bankEndpoint.String(), destinationAddress, proxyPrice, "3 FIL"},
		{bankEndpoint.String(), destinationAddress, proxyPrice, "2 FIL"},
		{bankEndpoint.String(), destinationAddress, proxyPrice, "1 FIL"},
		{bankEndpoint.String(), destinationAddress, proxyPrice, "0 FIL"},
		{bankEndpoint.String(), destinationAddress, proxyPrice, "wallet not found"},
	}

	for _, test := range tests {
		withdrawOpts := cli.WithdrawOptions{
			Amount:      test.amount,
			BankAddress: test.address,
			Destination: test.destination,
		}

		if err := cl.Validate.Struct(withdrawOpts); err != nil {
			t.Errorf("failed to validate: %v", err)
		}

		res, err := cli.Withdraw(ki, cfg.Wallet.Address, cfg.Route.Withdraw, withdrawOpts)
		if err != nil {
			if strings.Contains(err.Error(), test.expected) {
				continue
			}
			t.Errorf("failed to withdraw: %v", err)
		} else {
			assert.Equal(t, res.Status, "success")
			assert.Equal(t, res.Data.FIL.String(), test.expected)
		}
	}

	if err := setup.RunMigrations("DOWN", migr); err != nil {
		t.Fatalf("could not run down migrations: %v", err)
	}
}
