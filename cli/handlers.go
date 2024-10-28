package cli

import (
	"fmt"
	"io"
	"net/http"
)

func Deposit(bankAddress string) error {
	resp, err := PostRequest(bankAddress, "deposit", "application/json", "", 30)
	if err != nil {
		return fmt.Errorf("Error creating deposit request: %w", err)
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		_, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error reading deposit response body: %w", err)
		}
	case http.StatusNotFound:
		return fmt.Errorf("wallet not found")
	}

	return nil
}

func Withdraw(bankAddress string) error {
	resp, err := PostRequest(bankAddress, "withdraw", "application/json", "", 30)
	if err != nil {
		return fmt.Errorf("Error creating withdrawrequest: %w", err)
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		_, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error reading withdraw response body: %w", err)
		}
	case http.StatusNotFound:
		return fmt.Errorf("wallet not found")
	case http.StatusForbidden:
		return fmt.Errorf("insufficient funds")
	case http.StatusUnauthorized:
		return fmt.Errorf("funds locked for withdraw")
	}

	return nil
}
