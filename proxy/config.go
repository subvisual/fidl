package proxy

import (
	"log"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/subvisual/fidl/http"
	"github.com/subvisual/fidl/types"
)

type Bank struct {
	Register string `toml:"register"`
}

type Provider struct {
	Cost       types.FIL `toml:"cost"`
	SectorSize int64     `toml:"sector-size"`
}

type ForwarderConfig struct {
	DisableCompression bool          `toml:"disable-compression"`
	IdleConnTimeout    time.Duration `toml:"idle-conn-timeout"`
	HeaderTimeout      time.Duration `toml:"header-timeout"`
	MaxIdleConns       int           `toml:"max-idle-conns"`
	Upstream           string        `toml:"upstream"`
}

type Route struct {
	BankRedeem string `toml:"bank-redeem"`
	BankVerify string `toml:"bank-verify"`
}

type Config struct {
	Bank      map[string]Bank `toml:"bank"`
	Env       string          `toml:"env"`
	Forwarder ForwarderConfig `toml:"forwarder"`
	HTTP      http.HTTP       `toml:"http"`
	Logger    http.Logger     `toml:"logger"`
	Provider  Provider        `toml:"provider"`
	Route     Route           `toml:"route"`
	Wallet    types.Wallet    `toml:"wallet"`
}

func LoadConfiguration(cfgFilePath string) Config {
	var config Config
	if buf, err := os.ReadFile(cfgFilePath); err != nil {
		log.Fatalf("Config file not found: %s", cfgFilePath)
	} else if err := toml.Unmarshal(buf, &config); err != nil {
		log.Fatalf("Unable to parse configuration file: %v", err)
	}

	return config
}
