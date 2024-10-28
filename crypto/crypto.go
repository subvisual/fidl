package crypto

import (
	"fmt"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/venus/pkg/crypto"
	_ "github.com/filecoin-project/venus/pkg/crypto/secp" // to run init()
)

func Address(sigType crypto.SigType, pubkey []byte) (address.Address, error) {
	var addr address.Address
	var err error

	switch sigType {
	case crypto.SigTypeSecp256k1:
		addr, err = address.NewSecp256k1Address(pubkey)
	default:
		err = fmt.Errorf("no valid signature type: %v", sigType)
	}

	if err != nil {
		return address.Address{}, fmt.Errorf("failed to get address from pubkey: %w", err)
	}

	return addr, nil
}

func Verify(sig *crypto.Signature, addr address.Address, msg []byte) error {
	if err := crypto.Verify(sig, addr, msg); err != nil {
		return fmt.Errorf("failed to verify signature: %w", err)
	}

	return nil
}

func GeneratePrivate(sigType crypto.SigType) ([]byte, error) {
	pk, err := crypto.Generate(sigType)
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %w", err)
	}

	return pk, nil
}

func PrivateToPublic(sigType crypto.SigType, pk []byte) ([]byte, error) {
	pub, err := crypto.ToPublic(sigType, pk)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to public key: %w", err)
	}

	return pub, nil
}

func Sign(msg []byte, pk []byte, sigType crypto.SigType) (*crypto.Signature, error) {
	sig, err := crypto.Sign(msg, pk, sigType)
	if err != nil {
		return nil, fmt.Errorf("failed to sign message: %w", err)
	}

	return sig, nil
}
