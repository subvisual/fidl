package cli

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/subvisual/fidl/types"
)

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

	resp, err := PostRequest(cfg, options.BankAddress, "deposit", depositJSON, 30)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		depositResponse := TransactionResponse{}
		err := json.NewDecoder(resp.Body).Decode(&depositResponse)
		if err != nil {
			return fmt.Errorf("error decoding the response body: %s", err)
		}
		log.Printf("Deposit successful, funds updated to: %s", &depositResponse.Data.FIL)

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

	resp, err := PostRequest(cfg, options.BankAddress, "withdraw", withdrawJSON, 30)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		withdrawResponse := TransactionResponse{}
		err := json.NewDecoder(resp.Body).Decode(&withdrawResponse)
		if err != nil {
			return fmt.Errorf("error decoding the response body: %s", err)
		}
		log.Printf("Withdraw successful, your current bank balance is: %s", &withdrawResponse.Data.FIL)
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
