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

func (s BankService) Redeem(address string, uuid uuid.UUID, amount types.FIL) (bank.RedeemModel, error) {
	var spBalance types.FIL
	var cliBalance types.FIL
	var cliEscrow types.FIL
	var excess types.FIL

	verifyAuthQuery :=
		`
		SELECT *
		FROM escrow
		WHERE uuid = $1 
		  AND proxy = $2
		  AND balance >= $3
		  AND created_at >= $4
		  AND status_id = $5
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

	deleteAuthQuery :=
		`
		DELETE FROM escrow WHERE uuid = $1
		`

	cliEscrowQuery :=
		`
		UPDATE balances
  			SET escrow = escrow - $2,
				updated_at = now() at time zone 'utc'
  			WHERE id = $1
  			AND escrow >= $2
  			RETURNING escrow
		`

	deleteBalanceEntryQuery :=
		`
		DELETE FROM balances WHERE id = $1
		`

	deleteAccountEntryQuery :=
		`
		DELETE FROM accounts WHERE id = $1
		`

	err := Transaction(s.db, func(tx fidl.Queryable) error {
		account, err := getAccountByAddress(address, tx)
		if err != nil {
			return fmt.Errorf("failed to fetch account: %w", err)
		}

		if account.Type != StorageProvider {
			return bank.ErrOperationNotAllowed
		}

		var ok bool
		excess.Int, ok = new(big.Int).SetString(string("0"), 10)
		if !ok {
			return fmt.Errorf("failed to init excess to zero")
		}

		cliBalance.Int, ok = new(big.Int).SetString(string("0"), 10)
		if !ok {
			return fmt.Errorf("failed to init cliBalance to zero")
		}

		cfgDeadline, err := time.ParseDuration(s.cfg.EscrowDeadline)
		if err != nil {
			return fmt.Errorf("failed to parse escrow deadline from config: %w", err)
		}

		var auth Authorization

		args := []any{uuid, address, amount.Int.String(), time.Now().UTC().Add(-cfgDeadline), AuthorizationLocked}
		if err := tx.Get(&auth, verifyAuthQuery, args...); err != nil {
			return bank.ErrAuthNotFound
		}

		args = []any{account.ID, amount.Int.String()}
		if err := tx.QueryRow(depositQuery, args...).Scan(&spBalance); err != nil {
			return fmt.Errorf("failed to deposit balance to sp: %w", err)
		}

		args = []any{s.cfg.EscrowAddress, s.cfg.WalletAddress, amount.Int.String(), TransactionCompleted}
		if _, err := tx.Exec(transactionQuery, args...); err != nil {
			return fmt.Errorf("failed to register transaction during sp deposit: %w", err)
		}

		if auth.Balance.Int.Cmp(amount.Int) == 1 {
			excess.Int.Sub(auth.Balance.Int, amount.Int)

			args = []any{auth.ID, excess.Int.String()}
			if err := tx.QueryRow(depositQuery, args...).Scan(&cliBalance); err != nil {
				return fmt.Errorf("failed to deposit balance to cli: %w", err)
			}

			args = []any{s.cfg.EscrowAddress, s.cfg.WalletAddress, excess.Int.String(), TransactionCompleted}
			if _, err := tx.Exec(transactionQuery, args...); err != nil {
				return fmt.Errorf("failed to register transaction during cli deposit: %w", err)
			}
		}

		if _, err := tx.Exec(deleteAuthQuery, uuid); err != nil {
			return fmt.Errorf("failed to delete authorization during redeem: %w", err)
		}

		args = []any{auth.ID, auth.Balance.Int.String()}
		if err := tx.QueryRow(cliEscrowQuery, args...).Scan(&cliEscrow); err != nil {
			return fmt.Errorf("failed to update cli escrow: %w", err)
		}

		if cliBalance.Sign() == 0 && cliEscrow.Sign() == 0 {
			if _, err := tx.Exec(deleteBalanceEntryQuery, auth.ID); err != nil {
				return fmt.Errorf("failed to delete cli balance entry during redeem: %w", err)
			}

			if _, err := tx.Exec(deleteAccountEntryQuery, auth.ID); err != nil {
				return fmt.Errorf("failed to delete cli account entry during redeem: %w", err)
			}
		}

		return nil
	})
	if err != nil {
		return bank.RedeemModel{}, err
	}

	return bank.RedeemModel{
		Excess: excess,
		SP:     spBalance,
		CLI:    cliBalance,
	}, nil
}
