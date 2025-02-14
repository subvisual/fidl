package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/subvisual/fidl/cli"
	"github.com/subvisual/fidl/types"
)

func newBalanceCommand(cl cli.CLI) *cobra.Command {
	opts := cli.BalanceOptions{}
	balanceCmd := &cobra.Command{
		Use:   "balance",
		Short: "To check the client's account balance at a specified bank.",
		Long:  `This command checks the client's account balance at a specified bank.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			err := cl.Validate.Struct(opts)
			if err != nil {
				return fmt.Errorf("%w", err)
			}

			cfgPath, _ := cmd.Flags().GetString("config")
			cfg := cli.LoadConfiguration(cfgPath)

			ki, err := types.ReadWallet(cfg.Wallet)
			if err != nil {
				return fmt.Errorf("failed to read wallet: %w", err)
			}

			_, err = cli.Balance(ki, cfg.Wallet.Address, cfg.Route.Balance, opts)
			if err != nil {
				return fmt.Errorf("%w", err)
			}

			return nil
		},
	}

	balanceCmd.Flags().StringVarP(&opts.BankAddress, "bank", "b", "", "The bank address")
	cobra.CheckErr(balanceCmd.MarkFlagRequired("bank"))

	return balanceCmd
}
