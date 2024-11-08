package cli

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	fcrypto "github.com/subvisual/fidl/crypto"
	"github.com/subvisual/fidl/types"
)

func PostRequest(wallet types.Wallet, bankAddress string, route string, body []byte, timeout int64) (*http.Response, error) {
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

	ki, err := types.ReadWallet(wallet)
	if err != nil {
		return nil, fmt.Errorf("failed to read wallet: %w", err)
	}

	sig, err := fcrypto.Sign(ki.PrivateKey, types.SigTypeSecp256k1, msg)
	if err != nil {
		return nil, fmt.Errorf("failed to sign: %w", err)
	}

	sigBytes, err := sig.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal signature: %w", err)
	}

	req.Header.Add("content-type", "application/json")
	req.Header.Add("sig", hex.EncodeToString(sigBytes))
	req.Header.Add("pub", wallet.Address.String())
	req.Header.Add("msg", hex.EncodeToString(msg))

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http POST request failed: %w", err)
	}

	return resp, nil
}

func GetRequest(wallet types.Wallet, bankAddress string, route string, timeout int64) (*http.Response, error) {
	endpoint, err := url.JoinPath(bankAddress, route)
	if err != nil {
		return nil, fmt.Errorf("Error joining endpoint path: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("Error creating request: %w", err)
	}

	msg := []byte(time.Now().UTC().String())

	ki, err := types.ReadWallet(wallet)
	if err != nil {
		return nil, fmt.Errorf("failed to read wallet: %w", err)
	}

	sig, err := fcrypto.Sign(ki.PrivateKey, types.SigTypeSecp256k1, msg)
	if err != nil {
		return nil, fmt.Errorf("failed to sign: %w", err)
	}

	sigBytes, err := sig.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal signature: %w", err)
	}

	req.Header.Add("content-type", "application/json")
	req.Header.Add("sig", hex.EncodeToString(sigBytes))
	req.Header.Add("pub", wallet.Address.String())
	req.Header.Add("msg", hex.EncodeToString(msg))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making GET request: %w", err)
	}

	return resp, nil
}
