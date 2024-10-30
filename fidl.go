package fidl

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

var (
	//nolint
	Version string
	//nolint
	Commit string
)

type Queryable interface {
	Get(dest interface{}, query string, args ...interface{}) error
	QueryRow(query string, args ...any) *sql.Row
	Select(dest interface{}, query string, args ...interface{}) error
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
	NamedExec(query string, arg interface{}) (sql.Result, error)
	Exec(query string, args ...any) (sql.Result, error)
	PrepareNamed(query string) (*sqlx.NamedStmt, error)
}
