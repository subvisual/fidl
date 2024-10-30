package request

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/subvisual/fidl/crypto"
	"github.com/subvisual/fidl/proxy"
)

func Register(cfg proxy.Config) error {
	body, err := json.Marshal(map[string]any{
		"id":    cfg.Wallet.Address,
		"price": cfg.Provider.Cost,
	})
	if err != nil {
		return fmt.Errorf("failed payload marshaling: %w", err)
	}

	sig, err := crypto.Sign(cfg.Wallet, body)
	if err != nil {
		return fmt.Errorf("failed to sign: %w", err)
	}

	sigBytes, err := sig.MarshalBinary()
	if err != nil {
		return fmt.Errorf("filed to marshal signature: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	buf := bytes.NewBuffer(body)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, cfg.Bank.Register, buf)
	if err != nil {
		return fmt.Errorf("failed context creation: %w", err)
	}

	req.Header.Add("content-type", "application/json")
	req.Header.Add("sig", hex.EncodeToString(sigBytes))
	req.Header.Add("pub", cfg.Wallet.Address.String())
	req.Header.Add("msg", hex.EncodeToString(body))

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("http request failed: %w", err)
	}

	defer resp.Body.Close()

	return nil
}
