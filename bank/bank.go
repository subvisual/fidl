package bank

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/subvisual/fidl/http"
	"github.com/subvisual/fidl/types"
)

var (
	ErrInsufficientFunds     = errors.New("insufficient funds")
	ErrTransactionNotAllowed = errors.New("transaction not allowed")
	ErrLockedFunds           = errors.New("locked funds")
)

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

type RegisterParams struct {
	ID    string    `validate:"required" json:"id"`
	Price types.FIL `validate:"required" json:"price"`
}

type DepositParams struct {
	Amount types.FIL `validate:"required" json:"amount"`
}

type WithdrawParams struct {
	Amount      types.FIL `validate:"required" json:"amount"`
	Destination string    `validate:"required" json:"dst"`
}

type AuthorizeParams struct {
	Amount types.FIL `validate:"required" json:"amount"`
}

type RedeemParams struct {
	/* TODO */
}

type Service interface {
	RegisterProxy(spid string, source string, price types.FIL) error
	Deposit(address string, price types.FIL) (types.FIL, error)
	Withdraw(address string, destination string, price types.FIL) (types.FIL, error)
	Balance(address string) (types.FIL, types.FIL, error)
	Authorize(address string, amount types.FIL) (uuid.UUID, types.FIL, types.FIL, error)
}
