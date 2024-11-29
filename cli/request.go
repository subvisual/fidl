package cli

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/subvisual/fidl/crypto"
	"github.com/subvisual/fidl/request"
	"github.com/subvisual/fidl/types"
)

func PostRequest(ki types.KeyInfo, addr types.Address, bankAddress string, route string, body []byte) (*request.Response, error) {
	msg := append([]byte(time.Now().UTC().String()), body...)

	sig, err := sign(ki, msg)
	if err != nil {
		return nil, fmt.Errorf("failed to sign: %w", err)
	}

	dstURL, err := joinPath(bankAddress, route, "")
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	buff := bytes.NewBuffer(body)
	resp, err := request.New().
		SetEndpoint(dstURL).
		SetBody(buff).
		AppendHeader("content-type", "application/json").
		AppendHeader("sig", hex.EncodeToString(sig)).
		AppendHeader("pub", addr.String()).
		AppendHeader("msg", hex.EncodeToString(msg)).
		Post(context.Background())
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return resp, nil
}

func GetRequest(ki types.KeyInfo, addr types.Address, bankAddress string, route string, body []byte) (*request.Response, error) {
	msg := append([]byte(time.Now().UTC().String()), body...)

	sig, err := sign(ki, msg)
	if err != nil {
		return nil, fmt.Errorf("failed to sign: %w", err)
	}

	dstURL, err := joinPath(bankAddress, route, "")
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	buff := bytes.NewBuffer(body)
	resp, err := request.New().
		SetEndpoint(dstURL).
		SetBody(buff).
		AppendHeader("content-type", "application/json").
		AppendHeader("sig", hex.EncodeToString(sig)).
		AppendHeader("pub", addr.String()).
		AppendHeader("msg", hex.EncodeToString(msg)).
		Get(context.Background())
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return resp, nil
}

func ProxyRetrieveRequest(proxyAddress string, options RetrievalOptions, route string) (*request.Response, error) {
	uuid, err := uuid.Parse(options.Authorization)
	if err != nil {
		return nil, fmt.Errorf("error parsing authorization string to uuid: %w", err)
	}

	dstURL, err := joinPath(proxyAddress, route, options.Piece)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	resp, err := request.New().
		SetEndpoint(dstURL).
		AppendURLQuery("authorization", uuid.String()).
		Get(context.Background())
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return resp, nil
}

func ProxyBanksRequest(proxyAddress string, route string) (*request.Response, error) {
	dstURL, err := joinPath(proxyAddress, route, "")
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	resp, err := request.New().
		SetEndpoint(dstURL).
		Get(context.Background())
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return resp, nil
}

func sign(ki types.KeyInfo, body []byte) ([]byte, error) {
	sig, err := crypto.Sign(ki.PrivateKey, ki.Type, body)
	if err != nil {
		return nil, fmt.Errorf("failed to sign: %w", err)
	}

	out, err := sig.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("filed to marshal signature: %w", err)
	}

	return out, nil
}

func joinPath(address string, endpoint string, piece string) (*url.URL, error) {
	endpoint, err := url.JoinPath(address, endpoint, piece)
	if err != nil {
		return nil, fmt.Errorf("error joining endpoint path: %w", err)
	}

	dstURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("bank url: %w", err)
	}

	return dstURL, nil
}
