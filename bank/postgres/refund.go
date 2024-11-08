package postgres

import (
	"fmt"

	"github.com/subvisual/fidl"
	"github.com/subvisual/fidl/bank"
	"github.com/subvisual/fidl/types"
)

func (s BankService) Refund(address string) (bank.RefundBalances, error) {
	var expiredSum types.FIL
	var balance types.FIL
	var escrow types.FIL

	expiredQuery :=
		`
		SELECT COALESCE(SUM(balance), 0) AS expired_value_sum
		FROM escrow
		WHERE id = $1
		AND created_at < now() at time zone 'utc' - $2::interval
		`

	deleteExpiredQuery :=
		`
		DELETE FROM escrow
		WHERE id = $1
		AND created_at < now() at time zone 'utc' - $2::interval
		`

	updateBalancesQuery :=
		`
		UPDATE balances
  			SET balance = balance + $2,
				escrow = escrow - $2,
				updated_at = now() at time zone 'utc'
  			WHERE id = $1
  			AND escrow >= $2
  			RETURNING balance, escrow
		`

	// nolint:goconst
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

		args := []any{account.ID, s.cfg.EscrowDeadline}
		if err := tx.QueryRow(expiredQuery, args...).Scan(&expiredSum); err != nil {
			return fmt.Errorf("failed to get expired balance: %w", err)
		}

		if expiredSum.Sign() == 0 {
			return bank.ErrNothingToRefund
		}

		if _, err := tx.Exec(deleteExpiredQuery, args...); err != nil {
			return fmt.Errorf("failed to delete expired authorizations: %w", err)
		}

		args = []any{account.ID, expiredSum.Int.String()}
		if err := tx.QueryRow(updateBalancesQuery, args...).Scan(&balance, &escrow); err != nil {
			return fmt.Errorf("failed to update balances: %w", err)
		}

		args = []any{s.cfg.EscrowAddress, s.cfg.WalletAddress, expiredSum.Int.String(), bank.TransactionCompleted}
		if _, err := tx.Exec(transactionQuery, args...); err != nil {
			return fmt.Errorf("failed to register transaction during refund: %w", err)
		}

		return nil
	})
	if err != nil {
		return bank.RefundBalances{}, err
	}

	return bank.RefundBalances{
		Expired:   expiredSum,
		Available: balance,
		Escrow:    escrow,
	}, nil
}