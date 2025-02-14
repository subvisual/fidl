package blockchain

import (
	"context"
	"fmt"

	"github.com/defiweb/go-eth/abi"
	ethtypes "github.com/defiweb/go-eth/types"
	"github.com/subvisual/fidl/types"
)

func (c Client) Transfer(ctx context.Context, to string, amount types.FIL) (string, error) {
	transfer := abi.MustParseMethod("transfer(address, uint256)(bool)")

	calldata := transfer.MustEncodeArgs(to, amount.Int)

	tx := ethtypes.NewTransaction().
		SetTo(ethtypes.MustAddressFromHex(to)).
		SetValue(amount.Int).
		SetInput(calldata)

	txHash, _, err := c.SendTransaction(ctx, tx)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %w", err)
	}

	return txHash.String(), nil
}
