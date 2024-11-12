package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/subvisual/fidl/cli"
)

func newRefundCommand(cfg cli.Config) *cobra.Command {
	opts := cli.RefundOptions{}
	refundCmd := &cobra.Command{
		Use:   "refund",
		Short: "To refund FIL from the client's bank escrow.",
		Long: `This command transfers all the client's funds in escrow in a given bank, to the client's balance. For example:
	`,
		RunE: func(_ *cobra.Command, _ []string) error {
			err := cli.Refund(cfg, opts)
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
