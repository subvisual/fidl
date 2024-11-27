package postgres

import (
	"fmt"
	"time"

	"github.com/subvisual/fidl"
	"github.com/subvisual/fidl/bank"
	"github.com/subvisual/fidl/types"
)

func (s BankService) Refund(address string) (bank.RefundModel, error) {
	var expiredSum types.FIL
	var balance types.FIL
	var escrow types.FIL

	expiredQuery :=
		`
		SELECT COALESCE(SUM(balance), 0) AS expired_value_sum
		FROM escrow
		WHERE id = $1
		AND created_at < $2
		`

	deleteExpiredQuery :=
		`
		DELETE FROM escrow
		WHERE id = $1
		AND created_at < $2
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

		cfgDeadline, err := time.ParseDuration(s.cfg.EscrowDeadline)
		if err != nil {
			return fmt.Errorf("failed to parse escrow deadline from config: %w", err)
		}

		args := []any{account.ID, time.Now().UTC().Add(-cfgDeadline)}
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

		args = []any{s.cfg.EscrowAddress, s.cfg.WalletAddress, expiredSum.Int.String(), TransactionCompleted}
		if _, err := tx.Exec(transactionQuery, args...); err != nil {
			return fmt.Errorf("failed to register transaction during refund: %w", err)
		}

		return nil
	})
	if err != nil {
		return bank.RefundModel{}, err
	}

	return bank.RefundModel{
		Expired:   expiredSum,
		Available: balance,
		Escrow:    escrow,
	}, nil
}
