package bank

import (
	"github.com/subvisual/fidl"
	"github.com/subvisual/fidl/http"
)

type Server struct {
	*http.Server

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
	Deposit(source string, price fidl.FIL) (fidl.FIL, error)
	Withdraw(source string, price fidl.FIL) (fidl.FIL, error)
	Balance(source string) (fidl.FIL, error)
}
