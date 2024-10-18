package postgres

import (
	"fmt"

	"github.com/subvisual/fidl"
)

type BankService struct {
	db *DB
}

func NewBankService(db *DB) *BankService {
	return &BankService{
		db: db,
	}
}

func (s BankService) RegisterProxy(spid string, walletAddress string, price fidl.FIL) error {
	query :=
		`
		INSERT INTO storage_provider (sp_id, wallet_address, price)
		VALUES ($1, $2, $3)
		`

	args := []any{spid, walletAddress, price}
	if _, err := s.db.Exec(query, args...); err != nil {
		return fmt.Errorf("failed to add storage provider entry: %w", err)
	}

	return nil
}

func (s BankService) Deposit(source string, amount fidl.FIL) (fidl.FIL, error) {
	return 5.0, nil
}

func (s BankService) Withdraw(source string, amount fidl.FIL) (fidl.FIL, error) {
	return 5.0, nil
}

func (s BankService) Balance(source string) (fidl.FIL, error) {
	return 5.0, nil
}
