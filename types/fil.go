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

	switch t := value.(type) { // nolint:varnamelen
	case []uint8:
		var bInt big.Int
		_, ok := bInt.SetString(string(value.([]uint8)), 10)
		if !ok {
			return fmt.Errorf("failed to load value to []uint8: %v", value)
		}
		b.Int = &bInt
	default:
		return fmt.Errorf("could not scan type %T into FIL", t)
	}

	return nil
}
