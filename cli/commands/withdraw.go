package commands

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/subvisual/fidl/cli"
)

func newWithdrawCommand(bankAddress string) *cobra.Command {
	opts := cli.WithdrawOptions{}
	withdrawCmd := &cobra.Command{
		Use:   "withdraw",
		Short: "To withdraw FIL from the client's bank account.",
		Long: `This command transfers a specified amount of FIL from the client's bank account to their 
	own wallet. For example:
	`,
		Run: func(_ *cobra.Command, _ []string) {
			err := cli.Withdraw(bankAddress)
			if err != nil {
				log.Fatalf("failed to handle the withdraw request: %v", err)
			}
		},
	}

	withdrawCmd.Flags().Float64VarP(&opts.Amount, "amount", "a", 0, "The amount of funds to transfer")
	withdrawCmd.Flags().StringVarP(&opts.Publickey, "publickey", "p", "", "The wallet's public key")
	withdrawCmd.Flags().StringVarP(&opts.Signature, "signature", "s", "", "The signature to validate the wallet's ownership")
	cobra.CheckErr(withdrawCmd.MarkFlagRequired("amount"))
	cobra.CheckErr(withdrawCmd.MarkFlagRequired("publickey"))
	cobra.CheckErr(withdrawCmd.MarkFlagRequired("signature"))

	return withdrawCmd
}
