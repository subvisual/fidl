package postgres

import "time"

type AccountType int8

const (
	StorageProvider AccountType = iota + 1
	Client
)

func (a AccountType) String() string {
	switch a {
	case StorageProvider:
		return "Storage Provider"
	case Client:
		return "Client"
	default:
		return "Unknown" // nolint:goconst
	}
}

type Account struct {
	ID        int64       `db:"id"`
	Address   string      `db:"wallet_address"`
	Type      AccountType `db:"account_type"`
	CreatedAt time.Time   `db:"created_at"`
	UpdatedAt time.Time   `db:"updated_at"`
}
