package bank

import (
	"errors"
	"time"

	"github.com/subvisual/fidl"
	"github.com/subvisual/fidl/http"
)

var (
	ErrInsufficientFunds     = errors.New("insufficient funds")
	ErrTransactionNotAllowed = errors.New("transaction not allowed")
	ErrLockedFunds           = errors.New("locked funds")
)

type BalanceStatus int8

const (
	BalanceAvailable BalanceStatus = iota + 1
	BalanceLocked
)

func (a BalanceStatus) String() string {
	switch a {
	case BalanceAvailable:
		return "Available"
	case BalanceLocked:
		return "Locked"
	default:
		return "Unknown" // nolint:goconst
	}
}

type TransactionStatus int8

const (
	TransactionPending TransactionStatus = iota + 1
	TransactionCompleted
)

func (a TransactionStatus) String() string {
	switch a {
	case TransactionPending:
		return "Pending"
	case TransactionCompleted:
		return "Completed"
	default:
		return "Unknown" // nolint:goconst
	}
}

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

type Server struct {
	HTTP *http.Server

	BankService Service
}

type RegisterProxyParams struct {
	SpID      string   `validate:"required"`
	Signature string   `validate:"required"`
	PublicKey string   `validate:"required"`
	Price     fidl.FIL `validate:"required,gt=0.0"`
}

type TransactionParams struct {
	Signature string   `validate:"required"`
	PublicKey string   `validate:"required"`
	Amount    fidl.FIL `validate:"required,gt=0.0"`
}

type BalanceParams struct {
	Signature string `validate:"required"`
	PublicKey string `validate:"required"`
}

type Service interface {
	RegisterProxy(spid string, source string, price fidl.FIL) error
	Deposit(address string, price fidl.FIL) (fidl.FIL, error)
	Withdraw(address string, price fidl.FIL) (fidl.FIL, error)
	Balance(address string) (fidl.FIL, error)
}
