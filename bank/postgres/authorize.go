package postgres

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/subvisual/fidl"
	"github.com/subvisual/fidl/bank"
	"github.com/subvisual/fidl/types"
)

func (s BankService) Authorize(address string, proxy string) (bank.AuthModel, error) {
	var balance types.FIL
	var escrow types.FIL
	var cost types.FIL
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

	spCostQuery :=
		`
		SELECT price FROM storage_providers
		WHERE id = $1
		`

	escrowQuery :=
		`
		INSERT INTO escrow (id, uuid, balance, proxy, status_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING uuid, balance
		`

	// nolint:goconst
	transactionQuery :=
		`
		INSERT INTO transactions (transaction_id, source, destination, value, status_id)
		VALUES ($1, $2, $3, $4, $5)
		`

	err := Transaction(s.db, func(tx fidl.Queryable) error {
		account, err := getAccountByAddress(address, tx)
		if err != nil {
			return fmt.Errorf("failed to fetch cli account: %w", err)
		}

		spAccount, err := getAccountByAddress(proxy, tx)
		if err != nil {
			return fmt.Errorf("failed to fetch sp account: %w", err)
		}

		if err := tx.QueryRow(spCostQuery, spAccount.ID).Scan(&cost); err != nil {
			return fmt.Errorf("failed to fetch sp price: %w", err)
		}

		balance, _, err = s.Balance(address)
		if err != nil {
			return err
		}

		if cost.Cmp(balance.Int) == 1 {
			return bank.ErrInsufficientFunds
		}

		args := []any{account.ID, cost.Int.String()}
		if err := tx.QueryRow(withdrawQuery, args...).Scan(&balance); err != nil {
			return fmt.Errorf("failed to execute withdraw balance: %w", err)
		}

		transactionID, err := uuid.NewV7()
		if err != nil {
			return fmt.Errorf("failed to generate v7 uuid: %w", err)
		}

		uuid, err := uuid.NewV7()
		if err != nil {
			return fmt.Errorf("failed to generate v7 uuid: %w", err)
		}

		args = []any{account.ID, uuid, cost.Int.String(), proxy, AuthorizationOpen}
		if err := tx.QueryRow(escrowQuery, args...).Scan(&id, &escrow); err != nil {
			return fmt.Errorf("failed to deposit to escrow: %w", err)
		}

		args = []any{transactionID.String(), s.cfg.WalletAddress, s.cfg.EscrowAddress, cost.Int.String(), TransactionCompleted}
		if _, err := tx.Exec(transactionQuery, args...); err != nil {
			return fmt.Errorf("failed to register transaction during authorize: %w", err)
		}

		return nil
	})
	if err != nil {
		return bank.AuthModel{}, err
	}

	return bank.AuthModel{
		UUID:      id,
		Available: balance,
		Escrow:    escrow,
	}, nil
}
