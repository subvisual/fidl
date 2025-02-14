package main

import (
	"log"

	"github.com/go-playground/validator/v10"
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

	cl := cli.NewCLI(validator.New())

	if err := cl.RegisterValidators(); err != nil {
		log.Fatalf("Failed to register validators: %v", err)
	}

	cmd := commands.Parse(cl)
	_ = cmd.Execute()
}
