package setup

import (
	"fmt"
	"time"

	"github.com/subvisual/fidl/http"
	"github.com/subvisual/fidl/proxy"
	"github.com/subvisual/fidl/types"
)

func Proxy(price string, localhost string, bankPort int, upstreamPort int) (proxy.Config, error) {
	idleConnTimeout, _ := time.ParseDuration("90s")
	headerTimeout, _ := time.ParseDuration("10s")
	proxyAddress, _ := types.NewAddressFromString("f15nrbm4j7ptwl5tdbivasxnqypceh3cb4htpwlta")

	var cost types.FIL
	_ = cost.UnmarshalText([]byte(price))

	bankAddress := fmt.Sprintf("http://%s:%d", localhost, bankPort)

	cfg := proxy.Config{
		Bank: map[string]proxy.Bank{
			"one": {Register: bankAddress + "/api/v1/register"},
		},
		Env: "testing",
		Forwarder: proxy.ForwarderConfig{
			DisableCompression: true,
			IdleConnTimeout:    idleConnTimeout,
			HeaderTimeout:      headerTimeout,
			MaxIdleConns:       100,
			Upstream:           fmt.Sprintf("http://%s:%d", localhost, upstreamPort),
		},
		HTTP: http.HTTP{
			Addr:            "127.0.0.1",
			Fqdn:            localhost,
			Port:            bankPort + 1,
			ListenPort:      bankPort + 1,
			ReadTimeout:     15,
			WriteTimeout:    15,
			ShutdownTimeout: 10,
			TLS:             false,
		},
		Logger: http.Logger{
			Path:  "logs/proxy.log",
			Level: "FATAL",
		},
		Provider: proxy.Provider{
			Cost:       cost,
			SectorSize: 34359738368,
		},
		Wallet: types.Wallet{
			Address: proxyAddress,
			Path:    "etc/proxy.key",
		},
		Route: proxy.Route{
			BankRedeem: "/api/v1/redeem",
			BankVerify: "/api/v1/verify",
		},
	}

	return cfg, nil
}
