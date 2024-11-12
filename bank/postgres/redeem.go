package postgres

import (
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
	"github.com/subvisual/fidl"
	"github.com/subvisual/fidl/bank"
	"github.com/subvisual/fidl/types"
)

func (s BankService) Redeem(address string, uuid uuid.UUID, amount types.FIL) (bank.RedeemResponse, error) {
	var spBalance types.FIL
	var cliBalance types.FIL
	var excess types.FIL

	verifyAuthQuery :=
		`
		SELECT *
		FROM escrow
		WHERE (uuid = $1 AND balance >= $2)
		AND created_at < $3
		`

	depositQuery :=
		`
		UPDATE balances SET
			balance = balance + $2,
			updated_at = now() at time zone 'utc'
		WHERE id = $1
		RETURNING balance
		`

	// nolint:goconst
	transactionQuery :=
		`
		INSERT INTO transactions (source, destination, value, status_id)
		VALUES ($1, $2, $3, $4)
		`

	deleteQuery :=
		`
		DELETE FROM escrow WHERE uuid = $1
		`

	err := Transaction(s.db, func(tx fidl.Queryable) error {
		var ok bool
		excess.Int, ok = new(big.Int).SetString(string("0"), 10)
		if !ok {
			return fmt.Errorf("failed to init excess to zero")
		}

		var auth bank.Authorization

		cfgDeadline, err := time.ParseDuration(s.cfg.EscrowDeadline)
		if err != nil {
			return fmt.Errorf("failed to parse escrow deadline from config: %w", err)
		}

		args := []any{uuid, amount, time.Now().UTC().Add(-cfgDeadline)}
		if err := tx.Get(&auth, verifyAuthQuery, args...); err != nil {
			return bank.ErrAuthNotFound
		}

		account, err := getAccountByAddress(address, tx)
		if err != nil {
			return fmt.Errorf("failed to fetch account: %w", err)
		}

		if account.Type == bank.StorageProvider {
			return bank.ErrOperationNotAllowed
		}

		args = []any{account.ID, amount.Int.String()}
		if err := tx.QueryRow(depositQuery, args...).Scan(&spBalance); err != nil {
			return fmt.Errorf("failed to deposit balance to sp: %w", err)
		}

		args = []any{s.cfg.EscrowAddress, s.cfg.WalletAddress, amount.Int.String(), bank.TransactionCompleted}
		if _, err := tx.Exec(transactionQuery, args...); err != nil {
			return fmt.Errorf("failed to register transaction during sp deposit: %w", err)
		}

		if auth.Balance.Int.Cmp(amount.Int) == 1 {
			excess.Int.Sub(auth.Balance.Int, amount.Int)

			args := []any{auth.ID, excess.Int.String()}
			if err := tx.QueryRow(depositQuery, args...).Scan(&cliBalance); err != nil {
				return fmt.Errorf("failed to deposit balance to cli: %w", err)
			}

			args = []any{s.cfg.EscrowAddress, s.cfg.WalletAddress, excess.Int.String(), bank.TransactionCompleted}
			if _, err := tx.Exec(transactionQuery, args...); err != nil {
				return fmt.Errorf("failed to register transaction during cli deposit: %w", err)
			}
		}

		if _, err := tx.Exec(deleteQuery, uuid); err != nil {
			return fmt.Errorf("failed to delete authorization during redeem: %w", err)
		}

		return nil
	})
	if err != nil {
		return bank.RedeemResponse{}, err
	}

	return bank.RedeemResponse{
		Excess: excess,
		SP:     spBalance,
		CLI:    cliBalance,
	}, nil
}
