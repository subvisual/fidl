package http

type Logger struct {
	Level string `toml:"level"`
	Path  string `toml:"path"`
}

type HTTP struct {
	Addr            string `toml:"address"`
	Fqdn            string `toml:"fqdn"`
	Port            int    `toml:"port"`
	ListenPort      int    `toml:"listen-port"`
	ReadTimeout     int    `toml:"read-timeout"`
	WriteTimeout    int    `toml:"write-timeout"`
	ShutdownTimeout int    `toml:"shutdown-timeout"`
	TLS             bool   `toml:"tls"`
}
