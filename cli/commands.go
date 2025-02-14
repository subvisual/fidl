package cli

import (
	"github.com/google/uuid"
	"github.com/subvisual/fidl/types"
)

type AuthorizeOptions struct {
	BankAddress  string `validate:"url" json:"bankAddress"`
	ProxyInput   string `validate:"is-filecoin-address" json:"proxy"`
	ProxyAddress types.Address
}

type WithdrawOptions struct {
	Amount      string `json:"amount"`
	Destination string `validate:"is-valid-address" json:"dst"`
	BankAddress string `validate:"url" json:"bankAddress"`
}

type DepositOptions struct {
	Amount            string `json:"amount"`
	BankAddress       string `validate:"url" json:"bankAddress"`
	BankWalletAddress string `validate:"is-valid-address" json:"bankWalletAddress"`
	FIL               types.FIL
	TransactionHash   string
}

type BalanceOptions struct {
	BankAddress string `validate:"url" json:"bankAddress"`
}

type RefundOptions struct {
	BankAddress string `validate:"url" json:"bankAddress"`
}

type RetrievalOptions struct {
	ProxyAddress  string `validate:"url" json:"proxyAddress"`
	Piece         string `json:"piece"`
	Authorization string `validate:"uuid" json:"authorization"`
}

type BanksOptions struct {
	ProxyAddress string `validate:"url" json:"proxyAddress"`
}

type DepositResponseData struct {
	FIL types.FIL `json:"fil"`
}

type DepositResponse struct {
	Status string              `json:"status"`
	Data   DepositResponseData `json:"data"`
}

type WithdrawResponseData struct {
	FIL  types.FIL `json:"fil"`
	Hash string    `json:"hash"`
}

type WithdrawResponse struct {
	Status string               `json:"status"`
	Data   WithdrawResponseData `json:"data"`
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
