package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/subvisual/fidl/types"
)

func Deposit(cfg Config, options DepositOptions) error {
	var b types.FIL // nolint:varnamelen
	err := b.UnmarshalJSON([]byte(options.Amount))
	if err != nil {
		return fmt.Errorf("error unmarshaling amount data: %w", err)
	}

	depositJSON, err := json.Marshal(DepositBody{Amount: b})
	if err != nil {
		return fmt.Errorf("error marshaling deposit data: %w", err)
	}

	resp, err := PostRequest(cfg, "deposit", depositJSON, 30)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error reading deposit response body: %w", err)
		}
		bodyString := string(bodyBytes)
		log.Printf("Deposit succesfful, funds updated to: %s", bodyString)

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
		return fmt.Errorf("error unmarshaling amount data: %w", err)
	}

	withdrawJSON, err := json.Marshal(WithdrawBody{Amount: b, Destination: options.Destination})
	if err != nil {
		return fmt.Errorf("error marshaling withdraw data: %w", err)
	}

	resp, err := PostRequest(cfg, "withdraw", withdrawJSON, 30)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error reading withdraw response body: %w", err)
		}
		bodyString := string(bodyBytes)
		log.Printf("Withdraw succesfful, your current bank balance is: %s", bodyString)
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
