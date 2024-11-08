package cli

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/subvisual/fidl/types"
)

func Balance(cfg Config, options BalanceOptions) error {
	resp, err := GetRequest(cfg.Wallet, options.BankAddress, cfg.Route.Balance, 30)
	if err != nil {
		return fmt.Errorf("error creating balance request: %w", err)
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		balanceResponse := BalanceResponse{}
		err := json.NewDecoder(resp.Body).Decode(&balanceResponse)
		if err != nil {
			return fmt.Errorf("error decoding the response body: %w", err)
		}
		fmt.Printf("Your current bank balance is: %s \nYour current funds on escrow are: %s\n", balanceResponse.Data.FIL, balanceResponse.Data.Escrow) // nolint:forbidigo
	case http.StatusNotFound:
		return fmt.Errorf("wallet not found")
	default:
		return fmt.Errorf("something went wrong: %s", http.StatusText(resp.StatusCode))
	}

	return nil
}

func Deposit(cfg Config, options DepositOptions) error {
	var b types.FIL // nolint:varnamelen
	err := b.UnmarshalJSON([]byte(options.Amount))
	if err != nil {
		return fmt.Errorf("error unmarshalling amount data: %w", err)
	}

	depositJSON, err := json.Marshal(DepositBody{Amount: b})
	if err != nil {
		return fmt.Errorf("error marshalling deposit data: %w", err)
	}

	resp, err := PostRequest(cfg.Wallet, options.BankAddress, cfg.Route.Deposit, depositJSON, 30)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		depositResponse := TransactionResponse{}
		err := json.NewDecoder(resp.Body).Decode(&depositResponse)
		if err != nil {
			return fmt.Errorf("error decoding the response body: %w", err)
		}
		fmt.Println("Deposit successful, funds updated to:", depositResponse.Data.FIL) // nolint:forbidigo

	case http.StatusNotFound:
		return fmt.Errorf("wallet not found")
	default:
		return fmt.Errorf("something went wrong: %s", http.StatusText(resp.StatusCode))
	}

	return nil
}

func Withdraw(cfg Config, options WithdrawOptions) error {
	var b types.FIL // nolint:varnamelen
	err := b.UnmarshalJSON([]byte(options.Amount))
	if err != nil {
		return fmt.Errorf("error unmarshalling amount data: %w", err)
	}

	withdrawJSON, err := json.Marshal(WithdrawBody{Amount: b, Destination: options.Destination})
	if err != nil {
		return fmt.Errorf("error marshalling withdraw data: %w", err)
	}

	resp, err := PostRequest(cfg.Wallet, options.BankAddress, cfg.Route.Withdraw, withdrawJSON, 30)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		withdrawResponse := TransactionResponse{}
		err := json.NewDecoder(resp.Body).Decode(&withdrawResponse)
		if err != nil {
			return fmt.Errorf("error decoding the response body: %w", err)
		}
		fmt.Println("Withdraw successful, your current bank balance is:", withdrawResponse.Data.FIL) // nolint:forbidigo
	case http.StatusNotFound:
		return fmt.Errorf("wallet not found")
	case http.StatusForbidden:
		return fmt.Errorf("insufficient funds")
	case http.StatusUnauthorized:
		return fmt.Errorf("funds locked for withdraw")
	default:
		return fmt.Errorf("something went wrong: %s", http.StatusText(resp.StatusCode))
	}

	return nil
}
