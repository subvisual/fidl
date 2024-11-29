package setup

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/subvisual/fidl/cli"
	"github.com/subvisual/fidl/types"
)

func CLI() (cli.Config, cli.CLI, types.KeyInfo, error) {
	cfgFilePath := "../etc/cli.ini.example"
	cfg := cli.LoadConfiguration(cfgFilePath)

	cfg.Wallet.Path = "../" + cfg.Wallet.Path

	cl := cli.CLI{Validate: validator.New()}

	if err := cli.RegisterValidators(cl); err != nil {
		return cli.Config{}, cli.CLI{}, types.KeyInfo{}, fmt.Errorf("failed to register validators: %w", err)
	}

	ki, err := types.ReadWallet(cfg.Wallet)
	if err != nil {
		return cli.Config{}, cli.CLI{}, types.KeyInfo{}, fmt.Errorf("failed to read wallet: %w", err)
	}

	return cfg, cl, ki, nil
}
