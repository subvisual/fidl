package types

import (
	"encoding/json"
	"fmt"

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
