package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/subvisual/fidl/cli"
	"github.com/subvisual/fidl/types"
)

func newDepositCommand() *cobra.Command {
	opts := cli.DepositOptions{}
	depositCmd := &cobra.Command{
		Use:   "deposit",
		Short: "To deposit FIL into the client's bank account.",
		Long: `This command transfers a specified amount of FIL from the client's wallet to their
	account in the bank's system, securely updating the client's balance within the service.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfgPath, _ := cmd.Flags().GetString("config")
			cfg := cli.LoadConfiguration(cfgPath)

			ki, err := types.ReadWallet(cfg.Wallet)
			if err != nil {
				return fmt.Errorf("failed to read wallet: %w", err)
			}

			_, err = cli.Deposit(ki, cfg.Wallet.Address, cfg.Route.Deposit, opts)
			if err != nil {
				return fmt.Errorf("%w", err)
			}

			return nil
		},
	}

	depositCmd.Flags().StringVarP(&opts.Amount, "amount", "a", "", "The amount of funds to transfer")
	depositCmd.Flags().StringVarP(&opts.BankAddress, "bank", "b", "", "The bank address")
	cobra.CheckErr(depositCmd.MarkFlagRequired("amount"))
	cobra.CheckErr(depositCmd.MarkFlagRequired("bank"))

	return depositCmd
}
