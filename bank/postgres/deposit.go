package postgres

import (
	"fmt"

	"github.com/subvisual/fidl"
	"github.com/subvisual/fidl/bank"
	"github.com/subvisual/fidl/types"
)

func (s BankService) Deposit(address string, amount types.FIL) (types.FIL, error) {
	var balance types.FIL

	insertAccountQuery :=
		`
		INSERT INTO accounts (wallet_address, account_type)
		VALUES ($1, $2)
		ON CONFLICT (wallet_address) DO NOTHING
		`

	depositQuery :=
		`
		INSERT INTO balances (id, balance)
		VALUES ($1, $2)
		ON CONFLICT (id)
		DO UPDATE
		SET balance = balances.balance + EXCLUDED.balance,
			updated_at = now() at time zone 'utc'
		RETURNING balance
		`

	// nolint:goconst
	transactionQuery :=
		`
		INSERT INTO transactions (source, destination, value, status_id)
		VALUES ($1, $2, $3, $4)
		`

	err := Transaction(s.db, func(tx fidl.Queryable) error {
		args := []any{address, bank.Client}
		if _, err := tx.Exec(insertAccountQuery, args...); err != nil {
			return fmt.Errorf("failed to add account entry: %w", err)
		}

		account, err := getAccountByAddress(address, tx)
		if err != nil {
			return fmt.Errorf("failed to fetch account: %w", err)
		}

		if account.Type == bank.StorageProvider {
			return bank.ErrTransactionNotAllowed
		}

		args = []any{account.ID, amount.Int.String()}
		if err := tx.QueryRow(depositQuery, args...).Scan(&balance); err != nil {
			return fmt.Errorf("failed to deposit balance: %w", err)
		}

		args = []any{address, s.cfg.WalletAddress, amount.Int.String(), bank.TransactionCompleted}
		if _, err := tx.Exec(transactionQuery, args...); err != nil {
			return fmt.Errorf("failed to register transaction during deposit: %w", err)
		}

		return nil
	})
	if err != nil {
		return types.FIL{}, err
	}

	return balance, nil
}
