package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/subvisual/fidl/blockchain"
	"github.com/subvisual/fidl/cli"
	"github.com/subvisual/fidl/types"
)

func newDepositCommand(cl cli.CLI) *cobra.Command {
	opts := cli.DepositOptions{}
	depositCmd := &cobra.Command{
		Use:   "deposit",
		Short: "To deposit FIL into the client's bank account.",
		Long: `This command transfers a specified amount of FIL from the client's wallet to their
	account in the bank's system, securely updating the client's balance within the service.`,
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

			var fil types.FIL
			err = fil.UnmarshalJSON([]byte(opts.Amount))
			if err != nil {
				return fmt.Errorf("error unmarshalling amount data: %w", err)
			}

			opts.FIL = fil

			ethAddr, _, err := types.ParseAddress(opts.BankWalletAddress)
			if err != nil {
				return fmt.Errorf("failed to parse bank wallet public address: %w", err)
			}

			opts.BankWalletAddress = ethAddr

			blockchainService, err := blockchain.NewService(&blockchain.Config{
				RPCURL:                      cfg.Blockchain.RPCURL,
				GasLimitMultiplier:          cfg.Blockchain.GasLimitMultiplier,
				GasPriceMultiplier:          cfg.Blockchain.GasPriceMultiplier,
				PriorityFeePerGasMultiplier: cfg.Blockchain.PriorityFeePerGasMultiplier,
			}, ki.PrivateKey, 0)
			if err != nil {
				return fmt.Errorf("failed to create blockchain service: %w", err)
			}

			ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
			defer cancel()

			hash, err := blockchainService.Transfer(ctx, opts.BankWalletAddress, opts.FIL)
			if err != nil {
				return fmt.Errorf("failed to transfer funds: %w", err)
			}

			opts.TransactionHash = hash

			fmt.Println("Transferring funds, transaction hash:", hash) // nolint:forbidigo

			_, err = cli.Deposit(ctx, ki, cfg.Wallet.Address, cfg.Route.Deposit, opts)
			if err != nil {
				return fmt.Errorf("%w", err)
			}

			return nil
		},
	}

	depositCmd.Flags().StringVarP(&opts.Amount, "amount", "a", "", "The amount of funds to transfer")
	depositCmd.Flags().StringVarP(&opts.BankAddress, "bank", "b", "", "The bank address")
	depositCmd.Flags().StringVarP(&opts.BankWalletAddress, "bank-pub", "p", "", "The bank wallet public address")
	cobra.CheckErr(depositCmd.MarkFlagRequired("amount"))
	cobra.CheckErr(depositCmd.MarkFlagRequired("bank"))
	cobra.CheckErr(depositCmd.MarkFlagRequired("bank-pub"))

	return depositCmd
}
