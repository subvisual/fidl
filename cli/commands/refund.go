package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/subvisual/fidl/cli"
	"github.com/subvisual/fidl/types"
)

func newRefundCommand() *cobra.Command {
	opts := cli.RefundOptions{}
	refundCmd := &cobra.Command{
		Use:   "refund",
		Short: "To refund FIL from the client's bank escrow.",
		Long:  `This command transfers all the client's funds in escrow in a given bank, to the client's balance.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfgPath, _ := cmd.Flags().GetString("config")
			cfg := cli.LoadConfiguration(cfgPath)

			ki, err := types.ReadWallet(cfg.Wallet)
			if err != nil {
				return fmt.Errorf("failed to read wallet: %w", err)
			}

			_, err = cli.Refund(ki, cfg.Wallet.Address, cfg.Route.Refund, opts)
			if err != nil {
				return fmt.Errorf("%w", err)
			}

			return nil
		},
	}

	refundCmd.Flags().StringVarP(&opts.BankAddress, "bank", "b", "", "The bank address")
	cobra.CheckErr(refundCmd.MarkFlagRequired("bank"))

	return refundCmd
}
