package blockchain

type Config struct {
	RPCURL                      string  `toml:"rpc-url"`
	GasLimitMultiplier          float64 `toml:"gas-limit-multiplier"`
	GasPriceMultiplier          float64 `toml:"gas-price-multiplier"`
	PriorityFeePerGasMultiplier float64 `toml:"priority-fee-per-gas-multiplier"`
	VerifyInterval              int     `toml:"verify-interval"`
}
