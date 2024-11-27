package postgres

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/subvisual/fidl"
	"github.com/subvisual/fidl/bank"
	"github.com/subvisual/fidl/types"
)

func (s BankService) Verify(address string, uuid uuid.UUID, amount types.FIL) error {
	getAuthQuery :=
		`
			SELECT *
			FROM escrow
			WHERE uuid = $1
			  AND proxy = $2
			  AND balance >= $3
			  AND created_at >= $4
			  AND status_id = $5
			`

	updateAuthQuery :=
		`
		UPDATE escrow
  			SET status_id = $5,
				updated_at = now() at time zone 'utc'
  			WHERE uuid = $1
		  	  AND proxy = $2
		  	  AND balance >= $3
		  	  AND created_at >= $4
			  AND status_id = 1
		`

	err := Transaction(s.db, func(tx fidl.Queryable) error {
		account, err := getAccountByAddress(address, tx)
		if err != nil {
			return fmt.Errorf("failed to fetch account: %w", err)
		}

		if account.Type != StorageProvider {
			return bank.ErrOperationNotAllowed
		}

		var auth Authorization

		cfgDeadline, err := time.ParseDuration(s.cfg.EscrowDeadline)
		if err != nil {
			return fmt.Errorf("failed to parse escrow deadline from config: %w", err)
		}

		args := []any{uuid, address, amount.Int.String(), time.Now().UTC().Add(-cfgDeadline), AuthorizationOpen}
		if err := tx.Get(&auth, getAuthQuery, args...); err != nil {
			return bank.ErrAuthNotFound
		}

		args = []any{uuid, address, amount.Int.String(), time.Now().UTC().Add(-cfgDeadline), AuthorizationLocked}
		if _, err := tx.Exec(updateAuthQuery, args...); err != nil {
			return fmt.Errorf("failed to update authorization status: %w", err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
