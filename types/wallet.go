package types

import (
	"bytes"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

type Wallet struct {
	Path    string  `toml:"path"`
	Address Address `toml:"address"`
}

func ReadWallet(wallet Wallet) (KeyInfo, error) {
	var keyInfo KeyInfo

	pkIn, err := os.ReadFile(wallet.Path)
	if err != nil {
		return KeyInfo{}, fmt.Errorf("failed to load private key: %w", err)
	}

	pkIn = bytes.TrimRight(pkIn, "\n")
	pkey, err := hexutil.Decode(fmt.Sprintf("0x%s", string(pkIn)))
	if err != nil {
		return KeyInfo{}, fmt.Errorf("failed to decode private key: %w", err)
	}

	keyInfo.PrivateKey = pkey

	return keyInfo, nil
}
