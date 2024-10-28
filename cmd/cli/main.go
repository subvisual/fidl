package main

import (
	"flag"

	"github.com/subvisual/fidl"
	"github.com/subvisual/fidl/cli"
	"github.com/subvisual/fidl/cli/commands"
)

// nolint
var (
	version string
	commit  string
)

func main() {
	fidl.Version = version
	fidl.Commit = commit

	var cfgFilePath string
	flag.StringVar(&cfgFilePath, "config", "etc/cli.ini", "path to configuration file")
	flag.Parse()

	cfg := cli.LoadConfiguration(cfgFilePath)

	c := commands.Parse(cfg.CLI.BankAddress)
	_ = c.Execute()
}
