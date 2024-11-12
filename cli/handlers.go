package cli

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/subvisual/fidl/types"
)

func Authorize(cfg Config, options AuthorizeOptions) error {
	var b types.FIL // nolint:varnamelen
	err := b.UnmarshalJSON([]byte(options.Amount))
	if err != nil {
		return fmt.Errorf("error unmarshalling amount data: %w", err)
	}

	depositJSON, err := json.Marshal(DepositBody{Amount: b})
	if err != nil {
		return fmt.Errorf("error marshalling deposit data: %w", err)
	}

	resp, err := PostRequest(cfg.Wallet, options.BankAddress, cfg.Route.Authorize, depositJSON, 30)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		authorizeResponse := AuthorizeResponse{}
		err := json.NewDecoder(resp.Body).Decode(&authorizeResponse)
		if err != nil {
			return fmt.Errorf("error decoding the response body: %w", err)
		}
		fmt.Printf("Authorize successful, your current funds on escrow are: %s \nYour current bank balance is: %s\nAuth id: %s\n", authorizeResponse.Data.Escrow, authorizeResponse.Data.FIL, authorizeResponse.Data.ID) // nolint:forbidigo

	case http.StatusNotFound:
		return fmt.Errorf("wallet not found")
	case http.StatusForbidden:
		return fmt.Errorf("not have enough funds")
	default:
		return fmt.Errorf("something went wrong: %s", http.StatusText(resp.StatusCode))
	}

	return nil
}

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

func Refund(cfg Config, options RefundOptions) error {
	resp, err := GetRequest(cfg.Wallet, options.BankAddress, cfg.Route.Refund, 30)
	if err != nil {
		return fmt.Errorf("error creating refund request: %w", err)
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		refundResponse := RefundResponse{}
		err := json.NewDecoder(resp.Body).Decode(&refundResponse)
		if err != nil {
			return fmt.Errorf("error decoding the response body: %w", err)
		}
		fmt.Printf("Funds moved from escrow to your balance: %s \nYour current bank balance is: %s \nYour current funds on escrow are: %s\n", refundResponse.Data.Expired, refundResponse.Data.FIL, refundResponse.Data.Escrow) // nolint:forbidigo
	case http.StatusNotFound:
		return fmt.Errorf("wallet not found")
	case http.StatusUnprocessableEntity:
		return fmt.Errorf("no expired funds in escrow")
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
