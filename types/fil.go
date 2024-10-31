package types

import (
	"fmt"
	"math/big"

	"github.com/filecoin-project/venus/venus-shared/actors/types"
)

type FIL struct {
	types.FIL
}

func (b *FIL) Scan(value interface{}) error {
	if value == nil {
		b = nil
	}

	v, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("could not scan type %T into FIL", value)
	}

	b.Int, ok = new(big.Int).SetString(string(v), 10)
	if !ok {
		return fmt.Errorf("failed to load value to []uint8: %v", value)
	}

	return nil
}
