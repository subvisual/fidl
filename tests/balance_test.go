package tests

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/subvisual/fidl/blockchain"
	"github.com/subvisual/fidl/cli"
	"github.com/subvisual/fidl/proxy"
	"github.com/subvisual/fidl/tests/setup"
	"github.com/subvisual/fidl/types"
)

func TestBalance(t *testing.T) { // nolint:paralleltest
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

	// nolint:goconst
	amount := "5 FIL"

	var fil types.FIL
	err = fil.UnmarshalJSON([]byte(amount))
	if err != nil {
		t.Fatalf("error unmarshalling amount data: %v", err)
	}

	bankEthAddr, _, err := types.ParseAddress(bankWalletAddress)
	if err != nil {
		t.Fatalf("failed to parse bank wallet public address: %v", err)
	}

	blockchainService, err := blockchain.NewService(&blockchain.Config{
		RPCURL:                      cfg.Blockchain.RPCURL,
		GasLimitMultiplier:          cfg.Blockchain.GasLimitMultiplier,
		GasPriceMultiplier:          cfg.Blockchain.GasPriceMultiplier,
		PriorityFeePerGasMultiplier: cfg.Blockchain.PriorityFeePerGasMultiplier,
	}, ki.PrivateKey, 0)
	if err != nil {
		t.Fatalf("failed to create blockchain service: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	hash, err := blockchainService.Transfer(ctx, bankEthAddr, fil)
	if err != nil {
		t.Fatalf("failed to transfer funds: %v", err)
	}

	t.Logf("Transferring funds, transaction hash: %s", hash)

	depositOpts := cli.DepositOptions{
		Amount:            amount,
		BankAddress:       bankEndpoint.String(),
		BankWalletAddress: bankWalletAddress,
		FIL:               fil,
		TransactionHash:   hash,
	}

	if err := cl.Validate.Struct(depositOpts); err != nil {
		t.Errorf("failed to validate: %v", err)
	}

	res, err := cli.Deposit(ctx, ki, cfg.Wallet.Address, cfg.Route.Deposit, depositOpts)
	if err != nil {
		t.Fatalf("failed to deposit: %v", err)
	}

	assert.Equal(t, res.Status, "success")
	assert.Equal(t, res.Data.FIL.String(), "5 FIL")

	var tests = []struct {
		bankaddress  string
		proxyinput   string
		proxyaddress types.Address
		expected     string
		authorized   string
		escrow       string
	}{
		{bankEndpoint.String(), proxyCfg.Wallet.Address.String(), proxyCfg.Wallet.Address, "4 FIL", proxyPrice, "1 FIL"},
		{bankEndpoint.String(), proxyCfg.Wallet.Address.String(), proxyCfg.Wallet.Address, "3 FIL", proxyPrice, "2 FIL"},
		{bankEndpoint.String(), proxyCfg.Wallet.Address.String(), proxyCfg.Wallet.Address, "2 FIL", proxyPrice, "3 FIL"},
		{bankEndpoint.String(), proxyCfg.Wallet.Address.String(), proxyCfg.Wallet.Address, "1 FIL", proxyPrice, "4 FIL"},
		{bankEndpoint.String(), proxyCfg.Wallet.Address.String(), proxyCfg.Wallet.Address, "0 FIL", proxyPrice, "5 FIL"},
		{bankEndpoint.String(), proxyCfg.Wallet.Address.String(), proxyCfg.Wallet.Address, "not have enough funds", "", ""},
	}

	for _, test := range tests {
		authorizeOpts := cli.AuthorizeOptions{
			BankAddress:  test.bankaddress,
			ProxyInput:   test.proxyinput,
			ProxyAddress: test.proxyaddress,
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
		} else {
			assert.Equal(t, res.Status, "success")
			assert.Equal(t, res.Data.FIL.String(), test.expected)
			assert.Equal(t, res.Data.Escrow.String(), test.authorized)

			balanceOpts := cli.BalanceOptions{
				BankAddress: test.bankaddress,
			}

			if err := cl.Validate.Struct(balanceOpts); err != nil {
				t.Errorf("failed to validate: %v", err)
			}

			res, err := cli.Balance(ki, cfg.Wallet.Address, cfg.Route.Balance, balanceOpts)
			if err != nil {
				t.Errorf("failed to get balance: %v", err)
			}

			assert.Equal(t, res.Status, "success")
			assert.Equal(t, res.Data.FIL.String(), test.expected)
			assert.Equal(t, res.Data.Escrow.String(), test.escrow)
		}
	}

	if err := setup.RunMigrations("DOWN", migr); err != nil {
		t.Fatalf("could not run down migrations: %v", err)
	}
}
