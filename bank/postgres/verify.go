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
	query :=
		`
		SELECT *
		FROM escrow
		WHERE (uuid = $1 AND balance >= $2)
		AND created_at < $3
		`

	err := Transaction(s.db, func(tx fidl.Queryable) error {
		account, err := getAccountByAddress(address, tx)
		if err != nil {
			return fmt.Errorf("failed to fetch account: %w", err)
		}

		if account.Type != bank.StorageProvider {
			return bank.ErrOperationNotAllowed
		}

		var auth bank.Authorization

		cfgDeadline, err := time.ParseDuration(s.cfg.EscrowDeadline)
		if err != nil {
			return fmt.Errorf("failed to parse escrow deadline from config: %w", err)
		}

		args := []any{uuid, amount, time.Now().UTC().Add(-cfgDeadline)}

		if err := tx.Get(&auth, query, args...); err != nil {
			return bank.ErrAuthNotFound
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
