package proxy

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/subvisual/fidl/http"
)

type Wallet struct {
}

type Config struct {
	Env    string      `toml:"env"`
	Logger http.Logger `toml:"logger"`
	HTTP   http.HTTP   `toml:"http"`
	Wallet Wallet      `toml:"wallet"`
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
