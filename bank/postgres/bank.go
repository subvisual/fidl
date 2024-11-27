package postgres

import (
	"fmt"

	"github.com/subvisual/fidl"
)

type BankConfig struct {
	WalletAddress  string
	EscrowAddress  string
	EscrowDeadline string
}

type BankService struct {
	db  *DB
	cfg *BankConfig
}

func NewBankService(db *DB, cfg *BankConfig) *BankService {
	return &BankService{
		db:  db,
		cfg: cfg,
	}
}

func getAccountByAddress(address string, tx fidl.Queryable) (*Account, error) {
	query :=
		`
		SELECT *
		FROM accounts
		WHERE wallet_address = $1
		`

	var account Account
	if err := tx.Get(&account, query, address); err != nil {
		return nil, fmt.Errorf("failed to fetch account by wallet address: %w", err)
	}

	return &account, nil
}
