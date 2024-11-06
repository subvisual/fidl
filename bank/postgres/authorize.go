package postgres

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/subvisual/fidl"
	"github.com/subvisual/fidl/bank"
	"github.com/subvisual/fidl/types"
)

func (s BankService) Authorize(address string, amount types.FIL) (uuid.UUID, types.FIL, types.FIL, error) {
	var balance types.FIL
	var escrow types.FIL
	var id uuid.UUID

	withdrawQuery :=
		`
		UPDATE balances
  			SET balance = balance - $2,
				escrow = escrow + $2,
				updated_at = now() at time zone 'utc'
  			WHERE id = $1
  			AND balance >= $2
  			RETURNING balance
		`

	escrowQuery :=
		`
		INSERT INTO escrow (id, uuid, balance)
		VALUES ($1, $2, $3)
		RETURNING uuid, balance
		`

	transactionQuery :=
		`
		INSERT INTO transactions (source, destination, value, status_id)
		VALUES ($1, $2, $3, $4)
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

		uuid, err := uuid.NewV7()
		if err != nil {
			return fmt.Errorf("failed to generate v7 uuid: %w", err)
		}

		args = []any{account.ID, uuid, amount.Int.String()}
		if err := tx.QueryRow(escrowQuery, args...).Scan(&id, &escrow); err != nil {
			return fmt.Errorf("failed to deposit to escrow: %w", err)
		}

		args = []any{s.cfg.WalletAddress, s.cfg.EscrowAddress, amount.Int.String(), bank.TransactionCompleted}
		if _, err := tx.Exec(transactionQuery, args...); err != nil {
			return fmt.Errorf("failed to register transaction during authorize: %w", err)
		}

		return nil
	})
	if err != nil {
		return uuid.UUID{}, types.FIL{}, types.FIL{}, err
	}

	return id, balance, escrow, nil
}
