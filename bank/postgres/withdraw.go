package postgres

import (
	"fmt"

	"github.com/subvisual/fidl"
	"github.com/subvisual/fidl/bank"
	"github.com/subvisual/fidl/types"
)

func (s BankService) Withdraw(address string, destination string, amount types.FIL) (types.FIL, error) {
	var balance types.FIL

	if destination == s.cfg.WalletAddress {
		return types.FIL{}, bank.ErrTransactionNotAllowed
	}

	// nolint:goconst
	transactionQuery :=
		`
		INSERT INTO transactions (source, destination, value, status_id)
		VALUES ($1, $2, $3, $4)
		`

	withdrawQuery :=
		`
		UPDATE balances
  			SET balance = balance - $2,
				updated_at = now() at time zone 'utc'
  			WHERE id = $1
  			AND balance >= $2
  			RETURNING balance
		`

	err := Transaction(s.db, func(tx fidl.Queryable) error {
		account, err := getAccountByAddress(address, tx)
		if err != nil {
			return fmt.Errorf("failed to fetch account: %w", err)
		}

		balance, _, err = s.Balance(address)
		if err != nil {
			return err
		}

		if amount.Cmp(balance.Int) == 1 {
			return bank.ErrInsufficientFunds
		}

		args := []any{account.ID, amount.Int.String()}
		if err := tx.QueryRow(withdrawQuery, args...).Scan(&balance); err != nil {
			return fmt.Errorf("failed to execute withdraw balance: %w", err)
		}

		args = []any{s.cfg.WalletAddress, destination, amount.Int.String(), bank.TransactionCompleted}
		if _, err := tx.Exec(transactionQuery, args...); err != nil {
			return fmt.Errorf("failed to register transaction during withdraw: %w", err)
		}

		return nil
	})
	if err != nil {
		return types.FIL{}, err
	}

	return balance, nil
}
