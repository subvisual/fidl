package bank

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/subvisual/fidl/blockchain"
	"github.com/subvisual/fidl/http"
	"github.com/subvisual/fidl/types"
)

type Db struct {
	Dsn          string `toml:"dsn"`
	MaxOpenConns int    `toml:"max-open-connections"`
	MaxIdleConns int    `toml:"max-idle-connections"`
	MaxIdleTime  string `toml:"max-idle-time"`
}

type Escrow struct {
	Address  types.Address `toml:"address"`
	Deadline string        `toml:"deadline"`
}

type Config struct {
	Env        string            `toml:"env"`
	Logger     http.Logger       `toml:"logger"`
	Db         Db                `toml:"database"`
	HTTP       http.HTTP         `toml:"http"`
	Wallet     types.Wallet      `toml:"wallet"`
	Escrow     Escrow            `toml:"escrow"`
	Blockchain blockchain.Config `toml:"blockchain"`
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
