package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"time"

	"github.com/subvisual/fidl"
	"github.com/subvisual/fidl/http"
	"github.com/subvisual/fidl/proxy"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
	flag.StringVar(&cfgFilePath, "config", "etc/fidl.config", "path to configuration file")
	flag.Parse()

	cfg := proxy.LoadConfiguration(cfgFilePath)

	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() { <-c; cancel() }()

	zapcfg := zap.NewProductionConfig()
	zapcfg.OutputPaths = []string{cfg.Logger.Path, "stderr"}
	zapcfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	zapcfg.EncoderConfig.EncodeTime = zapcore.TimeEncoder(func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.UTC().Format(time.RFC3339))
	})

	abs, _ := filepath.Abs(cfg.Logger.Path)
	err := os.MkdirAll(path.Dir(abs), 0750)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	logger, err := zapcfg.Build()
	if err != nil {
		log.Fatalf("Failed to build zap logger: %v", err)
	}

	zap.ReplaceGlobals(logger)

	httpServer := http.New(&http.Config{
		Addr:            cfg.Proxy.Addr,
		Fqdn:            cfg.Proxy.Fqdn,
		Port:            cfg.Proxy.Port,
		ListenPort:      cfg.Proxy.ListenPort,
		ReadTimeout:     cfg.Proxy.ReadTimeout,
		WriteTimeout:    cfg.Proxy.WriteTimeout,
		ShutdownTimeout: cfg.Proxy.ShutdownTimeout,

		TLS: cfg.Proxy.TLS,
		Env: cfg.Env,
	})

	httpServer.Log = logger

	httpServer.RegisterValidator()

	if err := httpServer.RunProxy(); err != nil {
		logger.Fatal("failed to start http server", zap.Error(err))
	}

	// nolint
	defer logger.Sync()

	logger.Info("Server started", zap.String("addr", cfg.Proxy.Addr), zap.Int("port", cfg.Proxy.ListenPort))

	<-ctx.Done()

	logger.Info("Terminating...")

	if err := httpServer.Close(); err != nil {
		logger.Fatal("Error closing server connections", zap.Error(err))
	}
}
