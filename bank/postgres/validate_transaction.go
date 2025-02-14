package postgres

import (
	"fmt"
)

func (s BankService) ValidateBlockchainTransaction(hash string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM transactions
			WHERE transaction_id = $1
		)
	`

	var exists bool
	err := s.db.QueryRow(query, hash).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to validate transaction: %w", err)
	}

	if exists {
		return false, fmt.Errorf("transaction already registered")
	}

	return true, nil
}
