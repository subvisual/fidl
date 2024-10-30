package types

type Wallet struct {
	Path    string  `toml:"path"`
	Address Address `toml:"address"`
}
