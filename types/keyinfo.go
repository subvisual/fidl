package types

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	"github.com/filecoin-project/go-state-types/crypto"
)

type (
	Signature = crypto.Signature
)

const (
	SigTypeUnknown   = crypto.SigTypeUnknown
	SigTypeSecp256k1 = crypto.SigTypeSecp256k1
	SigTypeBLS       = crypto.SigTypeBLS
	SigTypeDelegated = crypto.SigTypeDelegated
)

type KeyInfo struct {
	Type       crypto.SigType
	PrivateKey []byte
}

func (ki *KeyInfo) UnmarshalJSON(value []byte) error {
	type KeyInfo struct {
		Type       string
		PrivateKey []byte
	}

	var keyInfo KeyInfo

	err := json.Unmarshal(value, &keyInfo)
	if err != nil {
		return fmt.Errorf("failed to unmarshal type: %w", err)
	}

	secp, _ := SigTypeSecp256k1.Name()
	bls, _ := SigTypeBLS.Name()
	del, _ := SigTypeDelegated.Name()
	ki.PrivateKey = keyInfo.PrivateKey

	switch keyInfo.Type {
	case secp:
		ki.Type = SigTypeSecp256k1
	case bls:
		ki.Type = SigTypeBLS
	case del:
		ki.Type = SigTypeDelegated
	default:
		ki.Type = SigTypeUnknown
	}

	return nil
}

func ReadWallet(wallet Wallet) (KeyInfo, error) {
	var keyInfo KeyInfo

	pkIn, err := os.ReadFile(wallet.Path)
	if err != nil {
		return KeyInfo{}, fmt.Errorf("failed to load private key: %w", err)
	}

	pkIn = bytes.TrimRight(pkIn, "\n")
	pkOut := make([]byte, hex.DecodedLen(len(pkIn)))
	if _, err := hex.Decode(pkOut, pkIn); err != nil {
		return KeyInfo{}, fmt.Errorf("failed to decode private key: %w", err)
	}

	if err := json.Unmarshal(pkOut, &keyInfo); err != nil {
		return KeyInfo{}, fmt.Errorf("failed to unmarshal private key: %w", err)
	}

	return keyInfo, nil
}
