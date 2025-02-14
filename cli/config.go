package cli

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/subvisual/fidl/blockchain"
	"github.com/subvisual/fidl/types"
)

type Route struct {
	Balance   string `toml:"balance"`
	Banks     string `toml:"banks"`
	Deposit   string `toml:"deposit"`
	Withdraw  string `toml:"withdraw"`
	Authorize string `toml:"authorize"`
	Refund    string `toml:"refund"`
	Retrieval string `toml:"retrieval"`
}

type Config struct {
	Env        string            `toml:"env"`
	Route      Route             `toml:"route"`
	Wallet     types.Wallet      `toml:"wallet"`
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
