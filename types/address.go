package types

import (
	"fmt"

	"github.com/filecoin-project/go-address"
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

func (a *Address) UnmarshalText(value []byte) error {
	addr, err := address.NewFromString(string(value))
	a.Address = &addr

	return err // nolint:wrapcheck
}
