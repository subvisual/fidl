package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/subvisual/fidl/cli"
	"github.com/subvisual/fidl/types"
)

func newWithdrawCommand(cl cli.CLI) *cobra.Command {
	opts := cli.WithdrawOptions{}
	withdrawCmd := &cobra.Command{
		Use:   "withdraw",
		Short: "To withdraw FIL from the client's bank account.",
		Long:  `This command transfers a specified amount of FIL from the client's bank account to their own wallet.`,
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

			_, err = cli.Withdraw(ki, cfg.Wallet.Address, cfg.Route.Withdraw, opts)
			if err != nil {
				return fmt.Errorf("%w", err)
			}

			return nil
		},
	}

	withdrawCmd.Flags().StringVarP(&opts.Amount, "amount", "a", "", "The amount of funds to transfer")
	withdrawCmd.Flags().StringVarP(&opts.Destination, "destination", "d", "", "The destination wallet")
	withdrawCmd.Flags().StringVarP(&opts.BankAddress, "bank", "b", "", "The bank address")
	cobra.CheckErr(withdrawCmd.MarkFlagRequired("amount"))
	cobra.CheckErr(withdrawCmd.MarkFlagRequired("destination"))
	cobra.CheckErr(withdrawCmd.MarkFlagRequired("bank"))

	return withdrawCmd
}
