package cli

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/subvisual/fidl/types"
)

type Route struct {
	Balance  string `toml:"balance"`
	Deposit  string `toml:"deposit"`
	Withdraw string `toml:"withdraw"`
}

type Config struct {
	Env    string       `toml:"env"`
	Route  Route        `toml:"route"`
	Wallet types.Wallet `toml:"wallet"`
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
