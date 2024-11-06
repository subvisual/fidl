package postgres

import (
	"fmt"

	"github.com/subvisual/fidl"
	"github.com/subvisual/fidl/types"
)

func (s BankService) Balance(address string) (types.FIL, error) {
	var balance types.FIL

	query :=
		`
		SELECT balance FROM balances WHERE id = $1
		`

	err := Transaction(s.db, func(tx fidl.Queryable) error {
		account, err := getAccountByAddress(address, tx)
		if err != nil {
			return fmt.Errorf("failed to fetch account: %w", err)
		}

		if err := s.db.Get(&balance, query, account.ID); err != nil {
			return fmt.Errorf("failed to balance by account id: %w", err)
		}

		return nil
	})
	if err != nil {
		return types.FIL{}, err
	}

	return balance, nil
}
