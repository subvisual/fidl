package proxy

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/subvisual/fidl/crypto"
	"github.com/subvisual/fidl/request"
	"github.com/subvisual/fidl/types"
	"go.uber.org/zap"
)

func Register(cfg Config) error {
	body, err := json.Marshal(map[string]any{
		"id":    cfg.Wallet.Address,
		"price": cfg.Provider.Cost,
	})
	if err != nil {
		return fmt.Errorf("failed payload marshaling: %w", err)
	}

	sig, err := sign(cfg.Wallet, body)
	if err != nil {
		return err
	}

	for key, val := range cfg.Bank {
		go func() {
			endpoint, _ := url.Parse(val.URL)
			if err := register(endpoint.JoinPath(cfg.Route.BankRegister), cfg.Wallet.Address.String(), body, sig); err != nil {
				zap.L().Error("failed to register bank", zap.String("bank", key), zap.Error(err))
			} else {
				zap.L().Info("registered with bank", zap.String("bank", key))
			}
		}()
	}

	return nil
}

func register(endpoint *url.URL, wallet string, payload []byte, sig []byte) error {
	buff := bytes.NewBuffer(payload)
	resp, err := request.New().
		SetEndpoint(endpoint).
		SetBody(buff).
		AppendHeader("content-type", "application/json").
		AppendHeader("sig", hex.EncodeToString(sig)).
		AppendHeader("pub", wallet).
		AppendHeader("msg", hex.EncodeToString(payload)).
		Post(context.Background())
	if err != nil {
		return fmt.Errorf("register bank: %w", err)
	}

	if resp.Status != http.StatusOK {
		return &request.Error{
			Message: resp.Body,
			Status:  resp.Status,
		}
	}

	return nil
}

func Verify(ctx context.Context, banks map[string]Bank, route Route, wallet types.Wallet, id uuid.UUID, amount types.FIL) (*Bank, error) {
	body, err := json.Marshal(map[string]any{
		"id":     id,
		"amount": amount,
	})
	if err != nil {
		return nil, fmt.Errorf("failed payload marshaling: %w", err)
	}

	sig, err := sign(wallet, body)
	if err != nil {
		return nil, err
	}

	for key, val := range banks {
		zap.L().Debug("looking up authorization at", zap.String("bank", key))
		endpoint, _ := url.Parse(val.URL)
		err := verify(ctx, endpoint.JoinPath(route.BankVerify), wallet, sig, body)
		if err != nil {
			zap.L().Debug("no authorization found at", zap.String("bank", key))
			continue
		}

		zap.L().Debug("authorization found at", zap.String("bank", key))

		return &val, nil
	}

	return nil, &request.Error{
		Message: []byte("no authorization found"),
		Status:  http.StatusNotFound,
	}
}

func verify(ctx context.Context, endpoint *url.URL, wallet types.Wallet, sig []byte, body []byte) error {
	buff := bytes.NewBuffer(body)
	resp, err := request.New().
		SetEndpoint(endpoint).
		SetBody(buff).
		AppendHeader("content-type", "application/json").
		AppendHeader("sig", hex.EncodeToString(sig)).
		AppendHeader("pub", wallet.Address.String()).
		AppendHeader("msg", hex.EncodeToString(body)).
		Post(ctx)
	if err != nil {
		return fmt.Errorf("failed to verify: %w", err)
	}

	if resp.Status != http.StatusOK {
		return &request.Error{
			Message: resp.Body,
			Status:  resp.Status,
		}
	}

	return nil
}

func Redeem(ctx context.Context, endpoint *url.URL, wallet types.Wallet, id uuid.UUID, amount types.FIL) error {
	body, err := json.Marshal(map[string]any{
		"id":     id,
		"amount": amount,
	})
	if err != nil {
		return fmt.Errorf("failed payload marshaling: %w", err)
	}

	sig, err := sign(wallet, body)
	if err != nil {
		return err
	}

	buff := bytes.NewBuffer(body)
	resp, err := request.New().
		SetEndpoint(endpoint).
		SetBody(buff).
		AppendHeader("content-type", "application/json").
		AppendHeader("sig", hex.EncodeToString(sig)).
		AppendHeader("pub", wallet.Address.String()).
		AppendHeader("msg", hex.EncodeToString(body)).
		Post(ctx)
	if err != nil {
		return fmt.Errorf("failed to redeem %w", err)
	}

	if resp.Status != http.StatusOK {
		return &request.Error{
			Message: resp.Body,
			Status:  resp.Status,
		}
	}

	return nil
}

func sign(wallet types.Wallet, body []byte) ([]byte, error) {
	ki, err := types.ReadWallet(wallet)
	if err != nil {
		return nil, fmt.Errorf("failed to read wallet: %w", err)
	}

	sig, err := crypto.Sign(ki.PrivateKey, types.SigTypeSecp256k1, body)
	if err != nil {
		return nil, fmt.Errorf("failed to sign: %w", err)
	}

	out, err := sig.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("filed to marshal signature: %w", err)
	}

	return out, nil
}
