package postgres

import (
	"fmt"
	"math/big"

	"github.com/subvisual/fidl"
	"github.com/subvisual/fidl/bank"
	"github.com/subvisual/fidl/types"
)

func (s BankService) Withdraw(address string, destination string, amount types.FIL) (types.FIL, error) {
	var balance types.FIL

	if destination == s.cfg.WalletAddress {
		return types.FIL{}, bank.ErrTransactionNotAllowed
	}

	transactionQuery :=
		`
		INSERT INTO transactions (source, destination, value, status_id)
		VALUES ($1, $2, $3, $4)
		`

	deleteQuery :=
		`
		WITH cte1 AS (
    		DELETE FROM balances WHERE id = $1
		)
		DELETE FROM accounts WHERE id = $1
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

		balance, err = s.Balance(address)
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

		zero, ok := new(big.Int).SetString("0", 10)
		if !ok {
			return fmt.Errorf("failed to set zero big int")
		}

		if balance.Cmp(zero) == 0 && account.Type == bank.Client {
			if _, err := tx.Exec(deleteQuery, account.ID); err != nil {
				return fmt.Errorf("failed to delete rows during withdraw balance: %w", err)
			}
		}

		return nil
	})
	if err != nil {
		return types.FIL{}, err
	}

	return balance, nil
}
