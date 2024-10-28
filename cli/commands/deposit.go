package commands

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/subvisual/fidl/cli"
)

func newDepositCommand(bankAddress string) *cobra.Command {
	opts := cli.DepositOptions{}
	depositCmd := &cobra.Command{
		Use:   "deposit",
		Short: "To deposit FIL into the client's bank account.",
		Long: `This command transfers a specified amount of FIL from the client's wallet to their
	account in the bank's system, securely updating the client's balance within the service. For example:
	`,
		Run: func(_ *cobra.Command, _ []string) {
			err := cli.Deposit(bankAddress)
			if err != nil {
				log.Fatalf("failed to handle the deposit request: %v", err)
			}
		},
	}
	depositCmd.Flags().Float64VarP(&opts.Amount, "amount", "a", 0, "The amount of funds to transfer")
	depositCmd.Flags().StringVarP(&opts.Publickey, "publickey", "p", "", "The wallet's public key")
	depositCmd.Flags().StringVarP(&opts.Signature, "signature", "s", "", "The signature to validate the wallet's ownership")
	cobra.CheckErr(depositCmd.MarkFlagRequired("amount"))
	cobra.CheckErr(depositCmd.MarkFlagRequired("publickey"))
	cobra.CheckErr(depositCmd.MarkFlagRequired("signature"))

	return depositCmd
}
