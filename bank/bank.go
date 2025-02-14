package bank

import (
	"time"

	"github.com/google/uuid"
	"github.com/subvisual/fidl/blockchain"
	"github.com/subvisual/fidl/http"
	"github.com/subvisual/fidl/types"
)

type Server struct {
	*http.Server

	BankService       Service
	BlockChainService blockchain.Service

	CustomReadTimeout time.Duration
}

type RegisterParams struct {
	ID    string    `validate:"required" json:"id"`
	Price types.FIL `validate:"required,is-valid-fil" json:"price"`
}

type DepositParams struct {
	Amount          types.FIL `validate:"required,is-valid-fil" json:"amount"`
	TransactionHash string    `validate:"required" json:"hash"`
}

type WithdrawParams struct {
	Amount      types.FIL `validate:"required,is-valid-fil" json:"amount"`
	Destination string    `validate:"required,is-valid-address" json:"dst"`
}

type AuthorizeParams struct {
	Proxy string `validate:"required,is-filecoin-address" json:"proxy"`
}

type RedeemParams struct {
	UUID   uuid.UUID `validate:"required" json:"id"`
	Amount types.FIL `validate:"required,is-valid-fil" json:"amount"`
}

type VerifyParams struct {
	UUID   uuid.UUID `validate:"required" json:"id"`
	Amount types.FIL `validate:"required,is-valid-fil" json:"amount"`
}

type RefundModel struct {
	Available types.FIL
	Escrow    types.FIL
	Expired   types.FIL
}

type AuthModel struct {
	UUID      uuid.UUID
	Available types.FIL
	Escrow    types.FIL
}

type RedeemModel struct {
	Excess types.FIL
	SP     types.FIL
	CLI    types.FIL
}

type Service interface {
	RegisterProxy(spid string, source string, price types.FIL) error
	ValidateBlockchainTransaction(hash string) (bool, error)
	Deposit(address string, price types.FIL, transactionHash string) (types.FIL, error)
	CanWithdraw(address string, amount types.FIL) (bool, error)
	Withdraw(address string, destination string, amount types.FIL, transactionHash string) (types.FIL, error)
	Balance(address string) (types.FIL, types.FIL, error)
	Authorize(address string, proxy string) (AuthModel, error)
	Refund(address string) (RefundModel, error)
	Verify(address string, uuid uuid.UUID, amount types.FIL) error
	Redeem(address string, uuid uuid.UUID, amount types.FIL) (RedeemModel, error)
}
