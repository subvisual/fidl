package postgres

import (
	"fmt"
	"math/big"

	"github.com/subvisual/fidl"
	"github.com/subvisual/fidl/bank"
	"github.com/subvisual/fidl/types"
)

type BankConfig struct {
	WalletAddress string
}

type BankService struct {
	db  *DB
	cfg *BankConfig
}

func NewBankService(db *DB, cfg *BankConfig) *BankService {
	return &BankService{
		db:  db,
		cfg: cfg,
	}
}

func getAccountByAddress(address string, tx fidl.Queryable) (*bank.Account, error) {
	query := `
		SELECT *
		FROM accounts
		WHERE wallet_address = $1`

	var account bank.Account
	if err := tx.Get(&account, query, address); err != nil {
		return nil, fmt.Errorf("failed to fetch account by wallet address: %w", err)
	}

	return &account, nil
}

func (s BankService) RegisterProxy(spid string, walletAddress string, price types.FIL) error {
	var accountID int64

	accountQuery :=
		`
		INSERT INTO accounts (wallet_address, account_type)
		VALUES ($1, $2)
		RETURNING id
		`

	balancesQuery :=
		`
		INSERT INTO balances (id)
		VALUES ($1)
		`

	spQuery :=
		`
		INSERT INTO storage_providers (id, sp_id, price)
		VALUES ($1, $2, $3)
		`

	err := Transaction(s.db, func(tx fidl.Queryable) error {
		args := []any{walletAddress, bank.StorageProvider}
		if err := tx.QueryRow(accountQuery, args...).Scan(&accountID); err != nil {
			return fmt.Errorf("failed to add account entry: %w", err)
		}

		if _, err := tx.Exec(balancesQuery, accountID); err != nil {
			return fmt.Errorf("failed to add balances entry: %w", err)
		}

		args = []any{accountID, spid, price.Int.String()}
		if _, err := tx.Exec(spQuery, args...); err != nil {
			return fmt.Errorf("failed to add storage provider entry: %w", err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (s BankService) Deposit(address string, amount types.FIL) (*types.FIL, error) {
	var balance types.FIL

	insertAccountQuery :=
		`
		INSERT INTO accounts (wallet_address, account_type)
		VALUES ($1, $2)
		ON CONFLICT (wallet_address) DO NOTHING
		`

	clientQuery :=
		`
		INSERT INTO clients (id)
		VALUES ($1)
		ON CONFLICT (id) DO NOTHING
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

		if _, err := tx.Exec(clientQuery, account.ID); err != nil {
			return fmt.Errorf("failed to add client entry: %w", err)
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
		return nil, err
	}

	return &balance, nil
}

func (s BankService) Withdraw(address string, destination string, amount types.FIL) (*types.FIL, error) {
	var balance *types.FIL

	if destination == s.cfg.WalletAddress {
		return nil, bank.ErrTransactionNotAllowed
	}

	transactionQuery :=
		`
		INSERT INTO transactions (source, destination, value, status_id)
		VALUES ($1, $2, $3, $4)
		`

	deleteQuery :=
		`
		WITH cte1 AS (
    		DELETE FROM clients WHERE id = $1
		), cte2 AS (
    		DELETE FROM balances WHERE id = $1
		)
		DELETE FROM accounts WHERE id = $1
		`

	withdrawQuery :=
		`
		UPDATE balances
  			SET balance = balance - $2,
				updated_at = now() at time zone 'utc'
  			WHERE id = $1
  			AND balance >= $2
  			RETURNING balance
		`

	err := Transaction(s.db, func(tx fidl.Queryable) error {
		account, err := getAccountByAddress(address, tx)
		if err != nil {
			return fmt.Errorf("failed to fetch account: %w", err)
		}

		if account.Type == bank.Client {
			status, err := s.BalanceStatus(account.ID)
			if err != nil {
				return err
			}

			if status != bank.BalanceAvailable {
				return bank.ErrLockedFunds
			}
		}

		balance, err = s.Balance(address)
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

		args = []any{s.cfg.WalletAddress, destination, amount.Int.String(), bank.TransactionCompleted}
		if _, err := tx.Exec(transactionQuery, args...); err != nil {
			return fmt.Errorf("failed to register transaction during withdraw: %w", err)
		}

		zero, ok := new(big.Int).SetString("0", 10)
		if !ok {
			return fmt.Errorf("failed to set zero big int")
		}

		if balance.Cmp(zero) == 0 && account.Type == bank.Client {
			if _, err := tx.Exec(deleteQuery, account.ID); err != nil {
				return fmt.Errorf("failed to delete rows during withdraw balance: %w", err)
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return balance, nil
}

func (s BankService) Balance(address string) (*types.FIL, error) {
	var balance types.FIL

	query :=
		`
		SELECT balance FROM balances WHERE id = $1
		`

	err := Transaction(s.db, func(tx fidl.Queryable) error {
		account, err := getAccountByAddress(address, tx)
		if err != nil {
			return fmt.Errorf("failed to fetch account: %w", err)
		}

		if err := s.db.Get(&balance, query, account.ID); err != nil {
			return fmt.Errorf("failed to balance by account id: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &balance, nil
}

func (s BankService) BalanceStatus(id int64) (bank.BalanceStatus, error) {
	var status bank.BalanceStatus

	query := `SELECT status_id FROM clients WHERE id = $1`
	if err := s.db.Get(&status, query, id); err != nil {
		return -1.0, fmt.Errorf("failed to balance status by account id: %w", err)
	}

	return status, nil
}
