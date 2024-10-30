package crypto

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/venus/pkg/crypto"
	_ "github.com/filecoin-project/venus/pkg/crypto/secp" // to run init()
	"github.com/subvisual/fidl"
)

func Verify(sig *crypto.Signature, addr address.Address, msg []byte) error {
	if err := crypto.Verify(sig, addr, msg); err != nil {
		return fmt.Errorf("failed to verify signature: %w", err)
	}

	return nil
}

func Sign(wallet fidl.Wallet, msg []byte) (*crypto.Signature, error) {
	var keyInfo KeyInfo

	pkIn, err := os.ReadFile(wallet.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to load private key: %w", err)
	}

	pkIn = bytes.TrimRight(pkIn, "\n")
	pkOut := make([]byte, hex.DecodedLen(len(pkIn)))
	if _, err := hex.Decode(pkOut, pkIn); err != nil {
		return nil, fmt.Errorf("failed to decode private key: %w", err)
	}

	if err := json.Unmarshal(pkOut, &keyInfo); err != nil {
		return nil, fmt.Errorf("failed to convert private key: %w", err)
	}

	sig, err := crypto.Sign(msg, keyInfo.PrivateKey, keyInfo.Type)
	if err != nil {
		return nil, fmt.Errorf("failed to sign message: %w", err)
	}

	return sig, nil
}
