package setup

import (
	"fmt"

	"github.com/subvisual/fidl/proxy"
	"github.com/subvisual/fidl/types"
)

func Proxy(price string) (proxy.Config, error) {
	cfgFilePath := "../etc/proxy.ini.example"
	cfg := proxy.LoadConfiguration(cfgFilePath)

	var cost types.FIL
	if err := cost.UnmarshalText([]byte(price)); err != nil {
		return proxy.Config{}, fmt.Errorf("failed to unmarshal proxy price: %w", err)
	}

	cfg.Provider.Cost = cost
	cfg.Wallet.Path = "../" + cfg.Wallet.Path

	return cfg, nil
}
