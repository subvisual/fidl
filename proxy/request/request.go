package request

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/subvisual/fidl/crypto"
	"github.com/subvisual/fidl/proxy"
	"github.com/subvisual/fidl/types"
	"go.uber.org/zap"
)

func Register(cfg proxy.Config) error {
	body, err := json.Marshal(map[string]any{
		"id":    cfg.Wallet.Address,
		"price": cfg.Provider.Cost,
	})
	if err != nil {
		return fmt.Errorf("failed payload marshaling: %w", err)
	}

	ki, err := types.ReadWallet(cfg.Wallet)
	if err != nil {
		return fmt.Errorf("failed to read wallet: %w", err)
	}

	sig, err := crypto.Sign(ki.PrivateKey, types.SigTypeSecp256k1, body)
	if err != nil {
		return fmt.Errorf("failed to sign: %w", err)
	}

	sigBytes, err := sig.MarshalBinary()
	if err != nil {
		return fmt.Errorf("failed to marshal signature: %w", err)
	}

	for k, v := range cfg.Bank {
		go func() {
			if err := register(v.Register, cfg.Wallet.Address.String(), body, sigBytes); err != nil {
				zap.L().Error("failed to register bank", zap.String("bank", k), zap.Error(err))
			} else {
				zap.L().Info("registered with bank", zap.String("bank", k))
			}
		}()
	}

	return nil
}

func register(endpoint string, wallet string, payload []byte, sig []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	buf := bytes.NewBuffer(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, buf)
	if err != nil {
		return fmt.Errorf("failed context creation: %w", err)
	}

	req.Header.Add("content-type", "application/json")
	req.Header.Add("sig", hex.EncodeToString(sig))
	req.Header.Add("pub", wallet)
	req.Header.Add("msg", hex.EncodeToString(payload))

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("http request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read message body: %w", err)
		}

		return fmt.Errorf("failed with: %s : %s", resp.Status, body)
	}

	defer resp.Body.Close()

	return nil
}
