package cli

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/subvisual/fidl/crypto"
	"github.com/subvisual/fidl/types"
)

func PostRequest(ki types.KeyInfo, addr types.Address, bankAddress string, route string, body []byte, timeout int64) (*http.Response, error) {
	endpoint, err := url.JoinPath(bankAddress, route)
	if err != nil {
		return nil, fmt.Errorf("error joining endpoint path: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(string(body)))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	msg := append([]byte(time.Now().UTC().String()), body...)

	sig, err := crypto.Sign(ki.PrivateKey, ki.Type, msg)
	if err != nil {
		return nil, fmt.Errorf("failed to sign: %w", err)
	}

	sigBytes, err := sig.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal signature: %w", err)
	}

	req.Header.Add("content-type", "application/json")
	req.Header.Add("sig", hex.EncodeToString(sigBytes))
	req.Header.Add("pub", addr.String())
	req.Header.Add("msg", hex.EncodeToString(msg))

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http POST request failed: %w", err)
	}

	return resp, nil
}

func GetRequest(ki types.KeyInfo, addr types.Address, bankAddress string, route string, timeout int64) (*http.Response, error) {
	endpoint, err := url.JoinPath(bankAddress, route)
	if err != nil {
		return nil, fmt.Errorf("error joining endpoint path: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	msg := []byte(time.Now().UTC().String())

	sig, err := crypto.Sign(ki.PrivateKey, ki.Type, msg)
	if err != nil {
		return nil, fmt.Errorf("failed to sign: %w", err)
	}

	sigBytes, err := sig.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal signature: %w", err)
	}

	req.Header.Add("content-type", "application/json")
	req.Header.Add("sig", hex.EncodeToString(sigBytes))
	req.Header.Add("pub", addr.String())
	req.Header.Add("msg", hex.EncodeToString(msg))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making GET request: %w", err)
	}

	return resp, nil
}

func ProxyRetrieveRequest(proxyAddress string, options RetrievalOptions, route string, timeout int64) (*http.Response, error) {
	endpoint, err := url.JoinPath(proxyAddress, route, options.Piece)
	if err != nil {
		return nil, fmt.Errorf("error joining endpoint path: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	uuid, err := uuid.Parse(options.Authorization)
	if err != nil {
		return nil, fmt.Errorf("error parsing authorization string to uuid: %w", err)
	}

	q := req.URL.Query()
	q.Add("bank", options.BankAddress)
	q.Add("authorization", uuid.String())
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making GET request: %w", err)
	}

	return resp, nil
}

func ProxyBanksRequest(proxyAddress string, route string, timeout int64) (*http.Response, error) {
	endpoint, err := url.JoinPath(proxyAddress, route)
	if err != nil {
		return nil, fmt.Errorf("error joining endpoint path: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making GET request: %w", err)
	}

	return resp, nil
}
