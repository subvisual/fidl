package fidl

import (
	"database/sql"
	"fmt"
	"math/big"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/venus/venus-shared/actors/types"
	"github.com/jmoiron/sqlx"
)

var (
	//nolint
	Version string
	//nolint
	Commit string
)

type FIL struct {
	types.FIL
}

func (b *FIL) Scan(value interface{}) error {
	if value == nil {
		b = nil
	}

	switch t := value.(type) {
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

type Wallet struct {
	Path    string  `toml:"path"`
	Address Address `toml:"address"`
}

type Address struct {
	*address.Address
}

func (a *Address) UnmarshalText(value []byte) error {
	addr, err := address.NewFromString(string(value))
	a.Address = &addr

	return err // nolint:wrapcheck
}

type Queryable interface {
	Get(dest interface{}, query string, args ...interface{}) error
	QueryRow(query string, args ...any) *sql.Row
	Select(dest interface{}, query string, args ...interface{}) error
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
	NamedExec(query string, arg interface{}) (sql.Result, error)
	Exec(query string, args ...any) (sql.Result, error)
	PrepareNamed(query string) (*sqlx.NamedStmt, error)
}
