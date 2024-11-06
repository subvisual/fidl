package postgres

import (
	"fmt"

	"github.com/subvisual/fidl"
	"github.com/subvisual/fidl/types"
)

func (s BankService) Balance(address string) (types.FIL, types.FIL, error) {
	var balance types.FIL
	var escrow types.FIL

	query :=
		`
		SELECT balance, escrow FROM balances WHERE id = $1
		`

	err := Transaction(s.db, func(tx fidl.Queryable) error {
		account, err := getAccountByAddress(address, tx)
		if err != nil {
			return fmt.Errorf("failed to fetch account: %w", err)
		}

		args := []any{account.ID}
		if err := tx.QueryRow(query, args...).Scan(&balance, &escrow); err != nil {
			return fmt.Errorf("failed to get balances: %w", err)
		}

		return nil
	})
	if err != nil {
		return types.FIL{}, types.FIL{}, err
	}

	return balance, escrow, nil
}
