package blockchain

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/defiweb/go-eth/types"
	"github.com/subvisual/fidl/collections"
	ftypes "github.com/subvisual/fidl/types"
	"go.uber.org/zap"
)

type VerifyTransactionOptions struct {
	Hash  string
	From  string
	Value ftypes.FIL
}

func (c Client) VerifyTransaction(ctx context.Context, opts VerifyTransactionOptions) error {
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, c.verifyTimeout)
	defer cancel()

	var timeNow = time.Now()

	var hash types.Hash
	err := hash.UnmarshalText([]byte(opts.Hash))
	if err != nil {
		return fmt.Errorf("failed to unmarshal hash: %w", err)
	}

	zap.L().Debug("Verifying transaction", zap.String("hash", opts.Hash), zap.String("from", opts.From), zap.String("value", opts.Value.String()))
	for {
		receipt, err := c.GetTransactionReceipt(ctx, hash)
		if err != nil {
			return fmt.Errorf("failed to get transaction receipt: %w", err)
		}

		if receipt.Status == nil {
			select {
			case <-ctx.Done():
				return fmt.Errorf("transaction verification timeout: %w", ctx.Err())
			case <-time.After(c.verifyInterval):
				continue
			}
		}

		if *receipt.Status == 1 {
			timeElapsed := time.Since(timeNow)
			zap.L().Debug("Transaction completed", zap.String("hash", opts.Hash), zap.Duration("time_elapsed", timeElapsed))

			break
		}
	}

	tx, err := c.GetTransactionByHash(ctx, hash)
	if err != nil {
		return fmt.Errorf("failed to get transaction: %w", err)
	}

	if !ValidTransactionValue(tx.Value, opts.Value) {
		return fmt.Errorf("invalid transaction value")
	}

	if err := ValidTransactionFrom(tx.From.String(), opts.From); err != nil {
		return fmt.Errorf("invalid transaction 'from' address")
	}

	if err := c.ValidTransactionTo(ctx, tx.To.String()); err != nil {
		return fmt.Errorf("invalid transaction 'to' address")
	}

	if err := c.ValidTransactionTimestamp(ctx, tx.BlockHash); err != nil {
		return fmt.Errorf("invalid transaction timestamp")
	}

	zap.L().Debug("Transaction is valid", zap.String("hash", opts.Hash), zap.String("from", opts.From), zap.String("value", opts.Value.String()))

	return nil
}

func ValidTransactionValue(txValue *big.Int, optsValue ftypes.FIL) bool {
	return optsValue.Cmp(txValue) == 0
}

func ValidTransactionFrom(txFrom string, optsFrom string) error {
	fromEthAddr, _, err := ftypes.ParseAddress(optsFrom)
	if err != nil {
		return fmt.Errorf("failed to parse 'from' address: %w", err)
	}

	if optsFrom != "" && fromEthAddr != txFrom {
		return fmt.Errorf("invalid transaction 'from' address")
	}

	return nil
}

func (c Client) ValidTransactionTo(ctx context.Context, txTo string) error {
	bankAddresses, err := c.Accounts(ctx)
	if err != nil {
		return fmt.Errorf("failed to get bank addresses: %w", err)
	}

	if !collections.ContainsFn(bankAddresses, func(item types.Address) bool {
		return txTo == item.String()
	}) {
		return fmt.Errorf("invalid transaction 'to' address")
	}

	return nil
}

func (c Client) ValidTransactionTimestamp(ctx context.Context, blockHash *types.Hash) error {
	block, err := c.BlockByHash(ctx, *blockHash, false)
	if err != nil {
		return fmt.Errorf("failed to get block number: %w, for hash: %s", err, blockHash.String())
	}

	if time.Since(block.Timestamp) > c.verifyTimeout {
		return fmt.Errorf("transaction deadline exceeded")
	}

	return nil
}
