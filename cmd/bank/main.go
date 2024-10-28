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
	flag.StringVar(&cfgFilePath, "config", "etc/bank.ini", "path to configuration file")
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

	httpServer := http.New(&http.Config{
		Addr:            cfg.HTTP.Addr,
		Fqdn:            cfg.HTTP.Fqdn,
		Port:            cfg.HTTP.Port,
		ListenPort:      cfg.HTTP.ListenPort,
		ReadTimeout:     cfg.HTTP.ReadTimeout,
		WriteTimeout:    cfg.HTTP.WriteTimeout,
		ShutdownTimeout: cfg.HTTP.ShutdownTimeout,

		TLS: cfg.HTTP.TLS,
		Env: cfg.Env,
	})

	bankCtx := bank.Server{HTTP: httpServer}
	bankCtx.BankService = postgres.NewBankService(db, &postgres.BankConfig{
		WalletAddress: cfg.Wallet.Address,
	})
	bankCtx.RegisterValidators()

	httpServer.Log = logger
	httpServer.RegisterMiddleWare(bank.AuthorizationCtx())
	httpServer.RegisterRoutes(bankCtx.Routes)

	if err := httpServer.Run(); err != nil {
		logger.Fatal("failed to start http server", zap.Error(err))
	}

	// nolint
	defer logger.Sync()

	logger.Info("Server started", zap.String("addr", cfg.HTTP.Addr), zap.Int("port", cfg.HTTP.ListenPort))

	<-ctx.Done()

	logger.Info("Terminating...")

	if err := httpServer.Close(); err != nil {
		logger.Fatal("Error closing server connections", zap.Error(err))
	}
}
