package fidl

import (
	"database/sql"

	"github.com/filecoin-project/go-address"
	"github.com/jmoiron/sqlx"
)

var (
	//nolint
	Version string
	//nolint
	Commit string
)

type FIL float64

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
