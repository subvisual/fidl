package types

import (
	"github.com/filecoin-project/go-address"
)

type Address struct {
	*address.Address
}

func (a *Address) UnmarshalText(value []byte) error {
	addr, err := address.NewFromString(string(value))
	a.Address = &addr

	return err // nolint:wrapcheck
}
