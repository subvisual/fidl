package cli

import "github.com/subvisual/fidl/types"

type WithdrawOptions struct {
	Amount      string `json:"amount"`
	Destination string `json:"dst"`
}

type DepositOptions struct {
	Amount string
}

type WithdrawBody struct {
	Amount      types.FIL `json:"amount"`
	Destination string    `json:"dst"`
}

type DepositBody struct {
	Amount types.FIL
}
