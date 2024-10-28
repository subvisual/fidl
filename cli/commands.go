package cli

type WithdrawOptions struct {
	Amount    float64
	Publickey string
	Signature string
}

type DepositOptions struct {
	Amount    float64
	Publickey string
	Signature string
}
