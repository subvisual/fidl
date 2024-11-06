package postgres

import (
	"fmt"

	"github.com/subvisual/fidl"
	"github.com/subvisual/fidl/bank"
)

type BankConfig struct {
	WalletAddress string
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

func getAccountByAddress(address string, tx fidl.Queryable) (*bank.Account, error) {
	query := `
		SELECT *
		FROM accounts
		WHERE wallet_address = $1`

	var account bank.Account
	if err := tx.Get(&account, query, address); err != nil {
		return nil, fmt.Errorf("failed to fetch account by wallet address: %w", err)
	}

	return &account, nil
}
