package cli

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type CLI struct {
}

type Wallet struct {
	// TODO
}

type Config struct {
	Env    string `toml:"env"`
	CLI    CLI    `toml:"cli"`
	Wallet Wallet `toml:"wallet"`
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
