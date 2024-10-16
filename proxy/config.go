package proxy

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type Logger struct {
	Level string `toml:"level"`
	Path  string `toml:"proxy-path"`
}

type Proxy struct {
	Addr            string `toml:"address"`
	Fqdn            string `toml:"fqdn"`
	Port            int    `toml:"port"`
	ListenPort      int    `toml:"listen-port"`
	ReadTimeout     int    `toml:"read-timeout"`
	WriteTimeout    int    `toml:"write-timeout"`
	ShutdownTimeout int    `toml:"shutdown-timeout"`
	TLS             bool   `toml:"tls"`
}

type Wallet struct {
	// TODO
}

type Config struct {
	Env    string `toml:"env"`
	Logger Logger `toml:"logger"`
	Proxy  Proxy  `toml:"fidl-proxy"`
	Wallet Wallet `toml:"wallet-proxy"`
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
