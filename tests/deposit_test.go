package tests

import (
	"context"
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/subvisual/fidl/blockchain"
	"github.com/subvisual/fidl/cli"
	"github.com/subvisual/fidl/tests/setup"
	"github.com/subvisual/fidl/types"
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

	// nolint:goconst
	amount := "1 FIL"

	var tests = []struct {
		address  string
		amount   string
		expected string
	}{
		{bankEndpoint.String(), amount, "1 FIL"},
	}

	for _, test := range tests {
		var fil types.FIL
		err = fil.UnmarshalJSON([]byte(test.amount))
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
			Amount:            test.amount,
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
		assert.Equal(t, res.Data.FIL.String(), test.expected)
	}

	if err := setup.RunMigrations("DOWN", migr); err != nil {
		t.Fatalf("could not run down migrations: %v", err)
	}
}
