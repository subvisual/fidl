package blockchain

import (
	"context"
	"fmt"
	"time"

	"github.com/defiweb/go-eth/rpc"
	"github.com/defiweb/go-eth/rpc/transport"
	"github.com/defiweb/go-eth/txmodifier"
	"github.com/defiweb/go-eth/wallet"
	"github.com/subvisual/fidl/types"
)

type Client struct {
	*rpc.Client

	verifyTimeout  time.Duration
	verifyInterval time.Duration
}

func NewService(cfg *Config, pkey []byte, timeout time.Duration) (*Client, error) {
	key := wallet.NewKeyFromBytes(pkey)
	if key == nil {
		return nil, fmt.Errorf("failed to load private key")
	}

	transport, err := transport.NewHTTP(transport.HTTPOptions{URL: cfg.RPCURL})
	if err != nil {
		return nil, fmt.Errorf("failed to create transport: %w", err)
	}

	client, err := rpc.NewClient(
		rpc.WithTransport(transport),
		rpc.WithKeys(key),
		rpc.WithDefaultAddress(key.Address()),
		rpc.WithTXModifiers(
			txmodifier.NewGasLimitEstimator(txmodifier.GasLimitEstimatorOptions{
				Multiplier: cfg.GasLimitMultiplier,
			}),

			txmodifier.NewEIP1559GasFeeEstimator(txmodifier.EIP1559GasFeeEstimatorOptions{
				GasPriceMultiplier:          cfg.GasPriceMultiplier,
				PriorityFeePerGasMultiplier: cfg.PriorityFeePerGasMultiplier,
			}),

			txmodifier.NewNonceProvider(txmodifier.NonceProviderOptions{
				UsePendingBlock: true,
			}),

			txmodifier.NewChainIDProvider(txmodifier.ChainIDProviderOptions{
				Replace: false,
				Cache:   true,
			}),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	return &Client{
		Client:         client,
		verifyTimeout:  timeout,
		verifyInterval: time.Duration(cfg.VerifyInterval) * time.Second,
	}, nil
}

type Service interface {
	VerifyTransaction(ctx context.Context, opts VerifyTransactionOptions) error
	Transfer(ctx context.Context, to string, amount types.FIL) (string, error)
}
