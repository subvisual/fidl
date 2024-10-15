package bank

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type Logger struct {
	Level string `toml:"level"`
	Path  string `toml:"bank-path"`
}

type Db struct {
	Dsn          string `toml:"dsn"`
	MaxOpenConns int    `toml:"max-open-connections"`
	MaxIdleConns int    `toml:"max-idle-connections"`
	MaxIdleTime  string `toml:"max-idle-time"`
}

type Bank struct {
	Addr            string `toml:"address"`
	Fqdn            string `toml:"fqdn"`
	Port            int    `toml:"port"`
	ListenPort      int    `toml:"listen-port"`
	ReadTimeout     int    `toml:"read-timeout"`
	WriteTimeout    int    `toml:"write-timeout"`
	ShutdownTimeout int    `toml:"shutdown-timeout"`
	TLS             bool   `toml:"tls"`
}

type Config struct {
	Env    string `toml:"env"`
	Logger Logger `toml:"logger"`
	Db     Db     `toml:"database"`
	Bank   Bank   `toml:"fidl-bank"`
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
