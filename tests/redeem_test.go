package tests

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/subvisual/fidl/cli"
	"github.com/subvisual/fidl/proxy"
	"github.com/subvisual/fidl/tests/setup"
	"github.com/subvisual/fidl/types"
	"golang.org/x/exp/rand"
)

func TestRedeem(t *testing.T) { // nolint:paralleltest
	if err := setup.RunMigrations("UP", migr); err != nil {
		t.Fatalf("could not run up migrations: %v", err)
	}

	proxyCfg, err := setup.Proxy(proxyPrice)
	if err != nil {
		t.Fatalf("could not setup proxy info: %v", err)
	}

	if err := proxy.Register(proxyCfg); err != nil {
		t.Log("failed to register proxy", err)
		t.Fail()
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

	var tests = []struct {
		address     string
		destination string
		expected    string
		authorized  string
	}{
		{bankEndpoint.String(), proxyCfg.Wallet.Address.String(), "4 FIL", proxyPrice},
		{bankEndpoint.String(), proxyCfg.Wallet.Address.String(), "3 FIL", proxyPrice},
		{bankEndpoint.String(), proxyCfg.Wallet.Address.String(), "2 FIL", proxyPrice},
		{bankEndpoint.String(), proxyCfg.Wallet.Address.String(), "1 FIL", proxyPrice},
		{bankEndpoint.String(), proxyCfg.Wallet.Address.String(), "0 FIL", proxyPrice},
	}

	bankEndpoint.Path = proxyCfg.Route.BankRedeem

	for _, test := range tests {
		authorizeOpts := cli.AuthorizeOptions{
			BankAddress: test.address,
			Proxy:       test.destination,
		}

		if err := cl.Validate.Struct(authorizeOpts); err != nil {
			t.Errorf("failed to validate: %v", err)
		}

		res, err := cli.Authorize(ki, cfg.Wallet.Address, cfg.Route.Authorize, authorizeOpts)
		if err != nil {
			if strings.Contains(err.Error(), test.expected) {
				continue
			}
			t.Errorf("failed to authorize: %v", err)
		}

		assert.Equal(t, res.Status, "success")
		assert.Equal(t, res.Data.Escrow.String(), test.authorized)

		// nolint:gosec
		rand.Seed(uint64(time.Now().UnixNano()))
		randomFloat := rand.Float64()

		s := strconv.FormatFloat(randomFloat, 'f', -1, 64)

		var cost types.FIL
		_ = cost.UnmarshalText([]byte(s))

		t.Log("cost value: ", cost)

		ctx := context.Background()

		_, err = proxy.Verify(ctx, proxyCfg.Bank, proxyCfg.Route, proxyCfg.Wallet, res.Data.ID, cost)
		if err != nil {
			t.Errorf("failed to verify: %v", err)
		}

		if err := proxy.Redeem(ctx, &bankEndpoint, proxyCfg.Wallet, res.Data.ID, cost); err != nil {
			t.Errorf("failed to redeem: %v", err)
		}

		balanceOpts := cli.BalanceOptions{
			BankAddress: test.address,
		}

		if err := cl.Validate.Struct(balanceOpts); err != nil {
			t.Errorf("failed to validate: %v", err)
		}

		if _, err := cli.Balance(ki, cfg.Wallet.Address, cfg.Route.Balance, balanceOpts); err != nil {
			t.Errorf("failed to get balance: %v", err)
		}

		assert.Equal(t, res.Status, "success")
	}

	if err := setup.RunMigrations("DOWN", migr); err != nil {
		t.Fatalf("could not run down migrations: %v", err)
	}
}
