package types

import (
	"fmt"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/crypto"
	"github.com/filecoin-project/venus/venus-shared/actors/types"
)

type Address struct {
	*address.Address
}

func NewAddressFromString(value string) (Address, error) {
	addr, err := address.NewFromString(value)
	if err != nil {
		return Address{}, fmt.Errorf("failed to convert string to filecoin address: %w", err)
	}

	return Address{Address: &addr}, nil
}

func AddressProtocolToSigType(protocol byte) crypto.SigType {
	switch protocol {
	case 1:
		return crypto.SigTypeSecp256k1
	case 3:
		return crypto.SigTypeBLS
	case 4:
		return crypto.SigTypeDelegated
	default:
		return crypto.SigTypeUnknown
	}
}

func ParseAddress(addr string) (string, crypto.SigType, error) {
	ethAddr, err := types.ParseEthAddress(addr)
	if err == nil {
		return ethAddr.String(), 0, nil
	}

	filAddr, err := NewAddressFromString(addr)
	if err != nil {
		return "", 0, fmt.Errorf("invalid address format: %w", err)
	}

	ethAddr, err = types.EthAddressFromFilecoinAddress(*filAddr.Address)
	if err != nil {
		return "", 0, fmt.Errorf("invalid address format: %w", err)
	}

	return ethAddr.String(), AddressProtocolToSigType(filAddr.Protocol()), nil
}

func (a *Address) UnmarshalText(value []byte) error {
	addr, err := address.NewFromString(string(value))
	a.Address = &addr

	return err // nolint:wrapcheck
}
