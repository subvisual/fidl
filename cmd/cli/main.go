package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/subvisual/fidl"
	"github.com/subvisual/fidl/cli"
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

	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() { <-c; cancel() }()

	/*
		CLI stuff
	*/
	fmt.Println(cfg) // no err

	<-ctx.Done()
}
