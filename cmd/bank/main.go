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
	"github.com/subvisual/fidl/bank"
	"github.com/subvisual/fidl/http"
	"github.com/subvisual/fidl/postgres"
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

	cfg := bank.LoadConfiguration(cfgFilePath)

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

	db := postgres.Connect(postgres.Config{
		Dsn:          cfg.Db.Dsn,
		MaxOpenConns: cfg.Db.MaxOpenConns,
		MaxIdleConns: cfg.Db.MaxIdleConns,
		MaxIdleTime:  cfg.Db.MaxIdleTime,
	})

	httpBankServer := http.New(&http.Config{
		Addr:            cfg.Bank.Addr,
		Fqdn:            cfg.Bank.Fqdn,
		Port:            cfg.Bank.Port,
		ListenPort:      cfg.Bank.ListenPort,
		ReadTimeout:     cfg.Bank.ReadTimeout,
		WriteTimeout:    cfg.Bank.WriteTimeout,
		ShutdownTimeout: cfg.Bank.ShutdownTimeout,

		TLS: cfg.Bank.TLS,
		Env: cfg.Env,
	})

	httpBankServer.BankService = postgres.NewBankService(db)

	httpBankServer.Log = logger

	httpBankServer.RegisterValidator()

	if err := httpBankServer.RunBank(); err != nil {
		logger.Fatal("failed to start http server", zap.Error(err))
	}

	// nolint
	defer logger.Sync()

	logger.Info("Server started", zap.String("addr", cfg.Bank.Addr), zap.Int("port", cfg.Bank.ListenPort))

	<-ctx.Done()

	logger.Info("Terminating...")

	if err := httpBankServer.Close(); err != nil {
		logger.Fatal("Error closing server connections", zap.Error(err))
	}
}
