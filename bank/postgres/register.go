package postgres

import (
	"fmt"

	"github.com/subvisual/fidl"
	"github.com/subvisual/fidl/bank"
	"github.com/subvisual/fidl/types"
)

func (s BankService) RegisterProxy(spid string, walletAddress string, price types.FIL) error {
	var accountID int64

	accountQuery :=
		`
		INSERT INTO accounts (wallet_address, account_type)
		VALUES ($1, $2)
		RETURNING id
		`

	balancesQuery :=
		`
		INSERT INTO balances (id)
		VALUES ($1)
		`

	spQuery :=
		`
		INSERT INTO storage_providers (id, sp_id, price)
		VALUES ($1, $2, $3)
		`

	err := Transaction(s.db, func(tx fidl.Queryable) error {
		args := []any{walletAddress, bank.StorageProvider}
		if err := tx.QueryRow(accountQuery, args...).Scan(&accountID); err != nil {
			return fmt.Errorf("failed to add account entry: %w", err)
		}

		if _, err := tx.Exec(balancesQuery, accountID); err != nil {
			return fmt.Errorf("failed to add balances entry: %w", err)
		}

		args = []any{accountID, spid, price.Int.String()}
		if _, err := tx.Exec(spQuery, args...); err != nil {
			return fmt.Errorf("failed to add storage provider entry: %w", err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
