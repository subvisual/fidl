package tests

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/subvisual/fidl/cli"
	"github.com/subvisual/fidl/proxy"
	"github.com/subvisual/fidl/tests/setup"
)

func TestRefund(t *testing.T) { // nolint:paralleltest
	if err := setup.RunMigrations("UP", migr); err != nil {
		t.Fatalf("could not run up migrations: %v", err)
	}

	proxyCfg, err := setup.Proxy(proxyPrice, localhost, bankPort, upstreamPort)
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

	bankAddress := fmt.Sprintf("http://%s:%d", localhost, bankPort)

	depositOpts := cli.DepositOptions{
		Amount:      "5 FIL",
		BankAddress: bankAddress,
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
		{bankAddress, proxyCfg.Wallet.Address.String(), "4 FIL", proxyPrice},
		{bankAddress, proxyCfg.Wallet.Address.String(), "3 FIL", proxyPrice},
		{bankAddress, proxyCfg.Wallet.Address.String(), "2 FIL", proxyPrice},
		{bankAddress, proxyCfg.Wallet.Address.String(), "1 FIL", proxyPrice},
		{bankAddress, proxyCfg.Wallet.Address.String(), "0 FIL", proxyPrice},
		{bankAddress, proxyCfg.Wallet.Address.String(), "not have enough funds", ""},
	}

	for _, test := range tests {
		authorizeOpts := cli.AuthorizeOptions{
			BankAddress: test.address,
			Proxy:       test.destination,
		}

		res, err := cli.Authorize(ki, cfg.Wallet.Address, cfg.Route.Authorize, authorizeOpts, cl)
		if err != nil {
			if strings.Contains(err.Error(), test.expected) {
				continue
			}
			t.Errorf("failed to authorize: %v", err)
		} else {
			assert.Equal(t, res.Status, "success")
			assert.Equal(t, res.Data.FIL.String(), test.expected)
			assert.Equal(t, res.Data.Escrow.String(), test.authorized)

			query :=
				`
				UPDATE escrow
  					SET created_at = $2
  					WHERE uuid = $1
				`
			args := []any{res.Data.ID, time.Now().UTC().Add(-25 * time.Hour)}
			if _, err := db.Exec(query, args...); err != nil {
				t.Log("failed to updated created_at: ", err)
			}
		}
	}

	refundOpts := cli.RefundOptions{
		BankAddress: bankAddress,
	}

	refundRes, err := cli.Refund(ki, cfg.Wallet.Address, cfg.Route.Refund, refundOpts)
	if err != nil {
		t.Errorf("failed to refund: %v", err)
	}

	assert.Equal(t, refundRes.Status, "success")
	assert.Equal(t, refundRes.Data.FIL.String(), "5 FIL")
	assert.Equal(t, refundRes.Data.Expired.String(), "5 FIL")
	assert.Equal(t, refundRes.Data.Escrow.String(), "0 FIL")

	if err := setup.RunMigrations("DOWN", migr); err != nil {
		t.Fatalf("could not run down migrations: %v", err)
	}
}
