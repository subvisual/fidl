package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/subvisual/fidl/cli"
)

func newWithdrawCommand(cfg cli.Config) *cobra.Command {
	opts := cli.WithdrawOptions{}
	withdrawCmd := &cobra.Command{
		Use:   "withdraw",
		Short: "To withdraw FIL from the client's bank account.",
		Long: `This command transfers a specified amount of FIL from the client's bank account to their 
	own wallet. For example:
	`,
		RunE: func(_ *cobra.Command, _ []string) error {
			err := cli.Withdraw(cfg, opts)
			if err != nil {
				return fmt.Errorf("%w", err)
			}

			return nil
		},
	}

	withdrawCmd.Flags().StringVarP(&opts.Amount, "amount", "a", "", "The amount of funds to transfer")
	withdrawCmd.Flags().StringVarP(&opts.Destination, "destination", "d", "", "The destination wallet")
	cobra.CheckErr(withdrawCmd.MarkFlagRequired("amount"))
	cobra.CheckErr(withdrawCmd.MarkFlagRequired("destination"))

	return withdrawCmd
}
