package cli

import (
	"github.com/google/uuid"
	"github.com/subvisual/fidl/types"
)

type AuthorizeOptions struct {
	BankAddress string `json:"bankAddress"`
	Proxy       string `validate:"is-filecoin-address" json:"proxy"`
}

type WithdrawOptions struct {
	Amount      string `json:"amount"`
	Destination string `validate:"is-filecoin-address" json:"dst"`
	BankAddress string `json:"bankAddress"`
}

type DepositOptions struct {
	Amount      string `json:"amount"`
	BankAddress string `json:"bankAddress"`
}

type BalanceOptions struct {
	BankAddress string `json:"bankAddress"`
}

type RefundOptions struct {
	BankAddress string `json:"bankAddress"`
}

type RetrievalOptions struct {
	ProxyAddress  string `json:"proxyAddress"`
	Piece         string `json:"piece"`
	Authorization string `json:"authorization"`
}

type BanksOptions struct {
	ProxyAddress string `json:"proxyAddress"`
}

type TransactionResponseData struct {
	FIL types.FIL `json:"fil"`
}

type TransactionResponse struct {
	Status string                  `json:"status"`
	Data   TransactionResponseData `json:"data"`
}

type BalanceResponseData struct {
	FIL    types.FIL `json:"fil"`
	Escrow types.FIL `json:"escrow"`
}

type Bank struct {
	URL  string    `json:"url"`
	Cost types.FIL `json:"cost"`
}

type BalanceResponse struct {
	Status string              `json:"status"`
	Data   BalanceResponseData `json:"data"`
}

type BanksResponse struct {
	Status string `json:"status"`
	Data   []Bank `json:"data"`
}

type AuthorizeResponseData struct {
	FIL    types.FIL `json:"fil"`
	Escrow types.FIL `json:"escrow"`
	ID     uuid.UUID `json:"id"`
}

type AuthorizeResponse struct {
	Status string                `json:"status"`
	Data   AuthorizeResponseData `json:"data"`
}

type RefundResponseData struct {
	FIL     types.FIL `json:"fil"`
	Escrow  types.FIL `json:"escrow"`
	Expired types.FIL `json:"expired"`
}

type RefundResponse struct {
	Status string             `json:"status"`
	Data   RefundResponseData `json:"data"`
}
