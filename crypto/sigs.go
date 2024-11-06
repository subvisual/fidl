package crypto

import (
	"fmt"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/venus/pkg/crypto"
	_ "github.com/filecoin-project/venus/pkg/crypto/secp" // to run init()
)

func Verify(sig *crypto.Signature, addr address.Address, msg []byte) error {
	if err := crypto.Verify(sig, addr, msg); err != nil {
		return fmt.Errorf("failed to verify signature: %w", err)
	}

	return nil
}

func Sign(privkey []byte, sigType crypto.SigType, msg []byte) (*crypto.Signature, error) {
	sig, err := crypto.Sign(msg, privkey, sigType)
	if err != nil {
		return nil, fmt.Errorf("failed to sign message: %w", err)
	}

	return sig, nil
}
