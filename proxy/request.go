package proxy

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/subvisual/fidl/crypto"
	"github.com/subvisual/fidl/types"
	"go.uber.org/zap"
)

type Request struct {
	body     io.Reader
	endpoint url.URL
	headers  map[string]string
}

type Response struct {
	Body   []byte
	Status int
}

func NewRequest() *Request {
	return &Request{headers: make(map[string]string)}
}

func (r *Request) SetEndpoint(endpoint url.URL) *Request {
	r.endpoint = endpoint
	return r
}

func (r *Request) SetBody(body io.Reader) *Request {
	r.body = body
	return r
}

func (r *Request) AppendHeader(key string, value string) *Request {
	r.headers[key] = value
	return r
}

func (r *Request) Post(ctx context.Context) (*Response, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, r.endpoint.String(), r.body)
	if err != nil {
		return nil, fmt.Errorf("failed context creation: %w", err)
	}

	for k, v := range r.headers {
		req.Header.Add(k, v)
	}

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read message body: %w", err)
	}

	return &Response{Body: body, Status: resp.StatusCode}, nil
}

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

	for k, v := range cfg.Bank {
		go func() {
			if err := register(v.Register, cfg.Wallet.Address.String(), body, sig); err != nil {
				zap.L().Error("failed to register bank", zap.String("bank", k), zap.Error(err))
			} else {
				zap.L().Info("registered with bank", zap.String("bank", k))
			}
		}()
	}

	return nil
}

func register(endpoint string, wallet string, payload []byte, sig []byte) error {
	dstURL, err := url.Parse(endpoint)
	if err != nil {
		return fmt.Errorf("bank url: %w", err)
	}

	buff := bytes.NewBuffer(payload)
	resp, err := NewRequest().
		SetEndpoint(*dstURL).
		SetBody(buff).
		AppendHeader("content-type", "application/json").
		AppendHeader("sig", hex.EncodeToString(sig)).
		AppendHeader("pub", wallet).
		AppendHeader("msg", hex.EncodeToString(payload)).
		Post(context.Background())
	if err != nil {
		return err
	}

	if resp.Status != http.StatusOK {
		return fmt.Errorf("bank register: %s", resp.Body)
	}

	return nil
}

func Verify(ctx context.Context, endpoint url.URL, wallet types.Wallet, id uuid.UUID, amount types.FIL) error {
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
	resp, err := NewRequest().
		SetEndpoint(endpoint).
		SetBody(buff).
		AppendHeader("content-type", "application/json").
		AppendHeader("sig", hex.EncodeToString(sig)).
		AppendHeader("pub", wallet.Address.String()).
		AppendHeader("msg", hex.EncodeToString(body)).
		Post(ctx)
	if err != nil {
		return err
	}

	if resp.Status != http.StatusOK {
		return fmt.Errorf("bank register: %s", resp.Body)
	}

	return nil
}

func Redeem(ctx context.Context, endpoint url.URL, wallet types.Wallet, id uuid.UUID, amount types.FIL) error {
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
	resp, err := NewRequest().
		SetEndpoint(endpoint).
		SetBody(buff).
		AppendHeader("content-type", "application/json").
		AppendHeader("sig", hex.EncodeToString(sig)).
		AppendHeader("pub", wallet.Address.String()).
		AppendHeader("msg", hex.EncodeToString(body)).
		Post(ctx)
	if err != nil {
		return err
	}

	if resp.Status != http.StatusOK {
		return fmt.Errorf("bank register: %s", resp.Body)
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

func rebuildBankEndpoint(initialURL string, path string) (url.URL, error) {
	parsedEndpoint, err := url.Parse(initialURL)
	if err != nil {
		return url.URL{}, fmt.Errorf("failed to parse url: %w", err)
	}

	finalEndpoint := url.URL{
		Scheme: parsedEndpoint.Scheme,
		Host:   parsedEndpoint.Host,
		Path:   path,
	}

	return finalEndpoint, nil
}
