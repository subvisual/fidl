package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/subvisual/fidl/types"
)

func Authorize(ki types.KeyInfo, addr types.Address, route string, options AuthorizeOptions) (*AuthorizeResponse, error) {
	authorizeResponse := AuthorizeResponse{}

	body, err := json.Marshal(map[string]any{
		"proxy": options.ProxyAddress.String(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed payload marshaling: %w", err)
	}

	resp, err := PostRequest(context.Background(), ki, addr, options.BankAddress, route, body)
	if err != nil {
		return nil, err
	}

	switch resp.Status {
	case http.StatusOK:
		err := json.NewDecoder(bytes.NewReader(resp.Body)).Decode(&authorizeResponse)
		if err != nil {
			return nil, fmt.Errorf("error decoding the response body: %w", err)
		}
		fmt.Printf("You successfully authorized to escrow: %s \nYour current bank balance is: %s\nAuth id: %s\n", authorizeResponse.Data.Escrow, authorizeResponse.Data.FIL, authorizeResponse.Data.ID) // nolint:forbidigo
	case http.StatusNotFound:
		return nil, fmt.Errorf("client or proxy wallet not found")
	case http.StatusForbidden:
		return nil, fmt.Errorf("not have enough funds")
	case http.StatusUnauthorized:
		return nil, fmt.Errorf("the wallet address and signature do not match")
	default:
		return nil, fmt.Errorf("something went wrong: %s\nMessage: %s", http.StatusText(resp.Status), resp.Body)
	}

	return &authorizeResponse, nil
}

func Balance(ki types.KeyInfo, addr types.Address, route string, options BalanceOptions) (*BalanceResponse, error) {
	balanceResponse := BalanceResponse{}

	resp, err := GetRequest(ki, addr, options.BankAddress, route, nil)
	if err != nil {
		return nil, err
	}

	switch resp.Status {
	case http.StatusOK:
		err := json.NewDecoder(bytes.NewReader(resp.Body)).Decode(&balanceResponse)
		if err != nil {
			return nil, fmt.Errorf("error decoding the response body: %w", err)
		}
		fmt.Printf("Your current bank balance is: %s \nYour current funds on escrow are: %s\n", balanceResponse.Data.FIL, balanceResponse.Data.Escrow) // nolint:forbidigo
	case http.StatusNotFound:
		return nil, fmt.Errorf("wallet not found")
	case http.StatusUnauthorized:
		return nil, fmt.Errorf("the wallet address and signature do not match")
	default:
		return nil, fmt.Errorf("something went wrong: %s\nMessage: %s", http.StatusText(resp.Status), resp.Body)
	}

	return &balanceResponse, nil
}

func Banks(route string, options BanksOptions) (*BanksResponse, error) {
	banksResponse := BanksResponse{}

	resp, err := ProxyBanksRequest(options.ProxyAddress, route)
	if err != nil {
		return nil, fmt.Errorf("error creating banks request: %w", err)
	}

	switch resp.Status {
	case http.StatusOK:
		err := json.NewDecoder(bytes.NewReader(resp.Body)).Decode(&banksResponse)
		if err != nil {
			return nil, fmt.Errorf("error decoding the response body: %w", err)
		}
		fmt.Println("The proxy is registered in the following banks:") // nolint:forbidigo
		for _, b := range banksResponse.Data {
			fmt.Printf("Bank address: %s, with cost: %s\n", b.URL, b.Cost) // nolint:forbidigo
		}
	default:
		return nil, fmt.Errorf("something went wrong: %s\nMessage: %s", http.StatusText(resp.Status), resp.Body)
	}

	return &banksResponse, nil
}

func Deposit(ctx context.Context, ki types.KeyInfo, addr types.Address, route string, options DepositOptions) (*DepositResponse, error) {
	depositResponse := DepositResponse{}

	body, err := json.Marshal(map[string]any{
		"amount": options.FIL,
		"hash":   options.TransactionHash,
	})
	if err != nil {
		return nil, fmt.Errorf("failed payload marshaling: %w", err)
	}

	resp, err := PostRequest(ctx, ki, addr, options.BankAddress, route, body)
	if err != nil {
		return nil, err
	}

	switch resp.Status {
	case http.StatusOK:
		err := json.NewDecoder(bytes.NewReader(resp.Body)).Decode(&depositResponse)
		if err != nil {
			return nil, fmt.Errorf("error decoding the response body: %w", err)
		}
		fmt.Println("Deposit successful, funds updated to:", depositResponse.Data.FIL) // nolint:forbidigo
	case http.StatusNotFound:
		return nil, fmt.Errorf("wallet not found")
	case http.StatusUnauthorized:
		return nil, fmt.Errorf("the wallet address and signature do not match")
	case http.StatusConflict:
		return nil, fmt.Errorf("invalid transaction")
	default:
		return nil, fmt.Errorf("something went wrong: %s\nMessage: %s", http.StatusText(resp.Status), resp.Body)
	}

	return &depositResponse, nil
}

func Retrieval(route string, options RetrievalOptions) error {
	resp, err := ProxyRetrieveRequest(options.ProxyAddress, options, route)
	if err != nil {
		return fmt.Errorf("error creating retrieval request: %w", err)
	}

	switch resp.Status {
	case http.StatusOK:
		outputFile, err := os.Create("output.car")
		if err != nil {
			return fmt.Errorf("error creating output file: %w", err)
		}
		defer outputFile.Close()

		_, err = io.Copy(outputFile, bytes.NewReader(resp.Body))
		if err != nil {
			return fmt.Errorf("error saving piece to output.car: %w", err)
		}

		fmt.Println("Piece successfully saved as output.car") // nolint:forbidigo
	default:
		return fmt.Errorf("something went wrong: %s\nMessage: %s", http.StatusText(resp.Status), resp.Body)
	}

	return nil
}

func Refund(ki types.KeyInfo, addr types.Address, route string, options RefundOptions) (*RefundResponse, error) {
	refundResponse := RefundResponse{}

	resp, err := GetRequest(ki, addr, options.BankAddress, route, nil)
	if err != nil {
		return nil, err
	}

	switch resp.Status {
	case http.StatusOK:
		err := json.NewDecoder(bytes.NewReader(resp.Body)).Decode(&refundResponse)
		if err != nil {
			return nil, fmt.Errorf("error decoding the response body: %w", err)
		}
		fmt.Printf("Funds moved from escrow to your balance: %s \nYour current bank balance is: %s \nYour current funds on escrow are: %s\n", refundResponse.Data.Expired, refundResponse.Data.FIL, refundResponse.Data.Escrow) // nolint:forbidigo
	case http.StatusNotFound:
		return nil, fmt.Errorf("wallet not found")
	case http.StatusUnprocessableEntity:
		return nil, fmt.Errorf("no expired funds in escrow")
	case http.StatusUnauthorized:
		return nil, fmt.Errorf("the wallet address and signature do not match")
	default:
		return nil, fmt.Errorf("something went wrong: %s\nMessage: %s", http.StatusText(resp.Status), resp.Body)
	}

	return &refundResponse, nil
}

func Withdraw(ki types.KeyInfo, addr types.Address, route string, options WithdrawOptions) (*WithdrawResponse, error) {
	var b types.FIL // nolint:varnamelen
	withdrawResponse := WithdrawResponse{}

	err := b.UnmarshalJSON([]byte(options.Amount))
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling amount data: %w", err)
	}

	body, err := json.Marshal(map[string]any{
		"amount": b,
		"dst":    options.Destination,
	})
	if err != nil {
		return nil, fmt.Errorf("failed payload marshaling: %w", err)
	}

	resp, err := PostRequest(context.Background(), ki, addr, options.BankAddress, route, body)
	if err != nil {
		return nil, err
	}

	switch resp.Status {
	case http.StatusOK:
		err := json.NewDecoder(bytes.NewReader(resp.Body)).Decode(&withdrawResponse)
		if err != nil {
			return nil, fmt.Errorf("error decoding the response body: %w", err)
		}
		fmt.Println("Withdraw successful, your current bank balance is:", withdrawResponse.Data.FIL) // nolint:forbidigo
		fmt.Println("Transaction hash is:", withdrawResponse.Data.Hash)                              // nolint:forbidigo
	case http.StatusNotFound:
		return nil, fmt.Errorf("wallet not found")
	case http.StatusForbidden:
		return nil, fmt.Errorf("insufficient funds")
	case http.StatusUnauthorized:
		return nil, fmt.Errorf("the wallet address and signature do not match")
	default:
		return nil, fmt.Errorf("something went wrong: %s\nMessage: %s", http.StatusText(resp.Status), resp.Body)
	}

	return &withdrawResponse, nil
}
