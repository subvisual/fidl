package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/subvisual/fidl/types"
)

func Authorize(ki types.KeyInfo, addr types.Address, route string, options AuthorizeOptions, cli CLI) (*AuthorizeResponse, error) {
	authorizeResponse := AuthorizeResponse{}

	params := AuthorizeBody{Proxy: options.Proxy}

	err := cli.Validate.Struct(params)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	authorizeJSON, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("error marshalling authorize data: %w", err)
	}

	resp, err := PostRequest(ki, addr, options.BankAddress, route, authorizeJSON, 30)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		err := json.NewDecoder(resp.Body).Decode(&authorizeResponse)
		if err != nil {
			return nil, fmt.Errorf("error decoding the response body: %w", err)
		}
		fmt.Printf("You successfully authorized to escrow: %s \nYour current bank balance is: %s\nAuth id: %s\n", authorizeResponse.Data.Escrow, authorizeResponse.Data.FIL, authorizeResponse.Data.ID) // nolint:forbidigo
	case http.StatusNotFound:
		return nil, fmt.Errorf("wallet not found")
	case http.StatusForbidden:
		return nil, fmt.Errorf("not have enough funds")
	case http.StatusUnauthorized:
		return nil, fmt.Errorf("the wallet address and signature do not match")
	default:
		return nil, fmt.Errorf("something went wrong: %s", http.StatusText(resp.StatusCode))
	}

	return &authorizeResponse, nil
}

func Balance(ki types.KeyInfo, addr types.Address, route string, options BalanceOptions) (*BalanceResponse, error) {
	balanceResponse := BalanceResponse{}
	resp, err := GetRequest(ki, addr, options.BankAddress, route, 30)
	if err != nil {
		return nil, fmt.Errorf("error creating balance request: %w", err)
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		err := json.NewDecoder(resp.Body).Decode(&balanceResponse)
		if err != nil {
			return nil, fmt.Errorf("error decoding the response body: %w", err)
		}
		fmt.Printf("Your current bank balance is: %s \nYour current funds on escrow are: %s\n", balanceResponse.Data.FIL, balanceResponse.Data.Escrow) // nolint:forbidigo
	case http.StatusNotFound:
		return nil, fmt.Errorf("wallet not found")
	case http.StatusUnauthorized:
		return nil, fmt.Errorf("the wallet address and signature do not match")
	default:
		return nil, fmt.Errorf("something went wrong: %s", http.StatusText(resp.StatusCode))
	}

	return &balanceResponse, nil
}

func Banks(route string, options BanksOptions) (*BanksResponse, error) {
	banksResponse := BanksResponse{}
	resp, err := ProxyBanksRequest(options.ProxyAddress, route, 30)
	if err != nil {
		return nil, fmt.Errorf("error creating banks request: %w", err)
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		err := json.NewDecoder(resp.Body).Decode(&banksResponse)
		if err != nil {
			return nil, fmt.Errorf("error decoding the response body: %w", err)
		}
		fmt.Println("The proxy is registered in the following banks:") // nolint:forbidigo
		for _, b := range banksResponse.Data {
			fmt.Printf("Bank address: %s, with cost: %s\n", b.URL, b.Cost) // nolint:forbidigo
		}
	default:
		return nil, fmt.Errorf("something went wrong: %s", http.StatusText(resp.StatusCode))
	}

	return &banksResponse, nil
}

func Deposit(ki types.KeyInfo, addr types.Address, route string, options DepositOptions) (*TransactionResponse, error) {
	var b types.FIL // nolint:varnamelen
	depositResponse := TransactionResponse{}

	err := b.UnmarshalJSON([]byte(options.Amount))
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling amount data: %w", err)
	}

	depositJSON, err := json.Marshal(DepositBody{Amount: b})
	if err != nil {
		return nil, fmt.Errorf("error marshalling deposit data: %w", err)
	}

	resp, err := PostRequest(ki, addr, options.BankAddress, route, depositJSON, 30)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		err := json.NewDecoder(resp.Body).Decode(&depositResponse)
		if err != nil {
			return nil, fmt.Errorf("error decoding the response body: %w", err)
		}
		fmt.Println("Deposit successful, funds updated to:", depositResponse.Data.FIL) // nolint:forbidigo

	case http.StatusNotFound:
		return nil, fmt.Errorf("wallet not found")
	case http.StatusUnauthorized:
		return nil, fmt.Errorf("the wallet address and signature do not match")
	default:
		return nil, fmt.Errorf("something went wrong: %s", http.StatusText(resp.StatusCode))
	}

	return &depositResponse, nil
}

func Retrieval(route string, options RetrievalOptions) error {
	resp, err := ProxyRetrieveRequest(options.ProxyAddress, options, route, 30)
	if err != nil {
		return fmt.Errorf("error creating retrieval request: %w", err)
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		outputFile, err := os.Create("output.car")
		if err != nil {
			return fmt.Errorf("error creating output file: %w", err)
		}
		defer outputFile.Close()

		_, err = io.Copy(outputFile, resp.Body)
		if err != nil {
			return fmt.Errorf("error saving piece to output.car: %w", err)
		}

		fmt.Println("Piece successfully saved as output.car") // nolint:forbidigo
	case http.StatusNotFound:
		return fmt.Errorf("piece not found")
	default:
		return fmt.Errorf("something went wrong: %s", http.StatusText(resp.StatusCode))
	}

	return nil
}

func Refund(ki types.KeyInfo, addr types.Address, route string, options RefundOptions) (*RefundResponse, error) {
	refundResponse := RefundResponse{}
	resp, err := GetRequest(ki, addr, options.BankAddress, route, 30)
	if err != nil {
		return nil, fmt.Errorf("error creating refund request: %w", err)
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		err := json.NewDecoder(resp.Body).Decode(&refundResponse)
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
		return nil, fmt.Errorf("something went wrong: %s", http.StatusText(resp.StatusCode))
	}

	return &refundResponse, nil
}

func Withdraw(ki types.KeyInfo, addr types.Address, route string, options WithdrawOptions, cli CLI) (*TransactionResponse, error) {
	var b types.FIL // nolint:varnamelen
	withdrawResponse := TransactionResponse{}

	err := b.UnmarshalJSON([]byte(options.Amount))
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling amount data: %w", err)
	}

	params := WithdrawBody{Amount: b, Destination: options.Destination}

	err = cli.Validate.Struct(params)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	withdrawJSON, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("error marshalling withdraw data: %w", err)
	}

	resp, err := PostRequest(ki, addr, options.BankAddress, route, withdrawJSON, 30)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		err := json.NewDecoder(resp.Body).Decode(&withdrawResponse)
		if err != nil {
			return nil, fmt.Errorf("error decoding the response body: %w", err)
		}
		fmt.Println("Withdraw successful, your current bank balance is:", withdrawResponse.Data.FIL) // nolint:forbidigo
	case http.StatusNotFound:
		return nil, fmt.Errorf("wallet not found")
	case http.StatusForbidden:
		return nil, fmt.Errorf("insufficient funds")
	case http.StatusUnauthorized:
		return nil, fmt.Errorf("the wallet address and signature do not match")
	default:
		return nil, fmt.Errorf("something went wrong: %s", http.StatusText(resp.StatusCode))
	}

	return &withdrawResponse, nil
}
