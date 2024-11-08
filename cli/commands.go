package cli

import "github.com/subvisual/fidl/types"

type WithdrawOptions struct {
	Amount      string `json:"amount"`
	Destination string `json:"dst"`
	BankAddress string `json:"bankAddress"`
}

type DepositOptions struct {
	Amount      string `json:"amount"`
	BankAddress string `json:"bankAddress"`
}

type BalanceOptions struct {
	BankAddress string `json:"bankAddress"`
}

type WithdrawBody struct {
	Amount      types.FIL `json:"amount"`
	Destination string    `json:"dst"`
}

type DepositBody struct {
	Amount types.FIL `json:"amount"`
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

type BalanceResponse struct {
	Status string              `json:"status"`
	Data   BalanceResponseData `json:"data"`
}
