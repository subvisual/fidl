package setup

import (
	"context"
	"log"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"time"

	"github.com/subvisual/fidl/bank"
	"github.com/subvisual/fidl/bank/postgres"
	"github.com/subvisual/fidl/http"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Server(cfg bank.Config, db *postgres.DB) {
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() { <-c; cancel() }()

	zapcfg := zap.NewProductionConfig()
	zapcfg.OutputPaths = []string{cfg.Logger.Path, "stderr"}

	var loggerLevel zapcore.Level
	err := loggerLevel.UnmarshalText([]byte(cfg.Logger.Level))
	if err != nil {
		log.Print(err, ", default to: ", zapcore.DebugLevel.String())
		loggerLevel = zapcore.DebugLevel
	}
	zapcfg.Level = zap.NewAtomicLevelAt(loggerLevel)

	zapcfg.EncoderConfig.EncodeTime = zapcore.TimeEncoder(func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.UTC().Format(time.RFC3339))
	})

	abs, _ := filepath.Abs(cfg.Logger.Path)
	err = os.MkdirAll(path.Dir(abs), 0750)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	logger, err := zapcfg.Build()
	if err != nil {
		log.Fatalf("Failed to build zap logger: %v", err)
	}

	zap.ReplaceGlobals(logger)

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

	bankCtx := bank.Server{Server: httpServer}
	bankCtx.BankService = postgres.NewBankService(db, &postgres.BankConfig{
		WalletAddress:  cfg.Wallet.Address.String(),
		EscrowAddress:  cfg.Escrow.Address.String(),
		EscrowDeadline: cfg.Escrow.Deadline,
	})
	bankCtx.RegisterValidators()

	httpServer.Log = logger
	httpServer.RegisterMiddleWare()
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
