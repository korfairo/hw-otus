package main

import (
	"context"
	"flag"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/korfairo/hw-otus/hw12_13_14_15_calendar/internal/app"
	"github.com/korfairo/hw-otus/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/korfairo/hw-otus/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/korfairo/hw-otus/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/korfairo/hw-otus/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFilePath string

func init() {
	flag.StringVar(&configFilePath, "config", "configs/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	cfg, err := NewConfig(configFilePath)
	if err != nil {
		fmt.Println("failed to init config:", err.Error())
		return
	}

	log, err := logger.New(logger.Config{
		Level: cfg.Logger.Level,
	})
	if err != nil {
		fmt.Println("failed to init logger:", err.Error())
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	var storage app.Storage
	if cfg.Storage.Memory {
		storage = memorystorage.New()
	} else {
		extStorage := sqlstorage.New(cfg.Storage.PostgresDSN)
		if err = extStorage.Connect(ctx); err != nil {
			log.
				WithError(err).
				WithField("DSN", cfg.Storage.PostgresDSN).
				Fatal("couldn't connect to PostgreSQL DB")
		}
		defer extStorage.Close()
		storage = extStorage
	}

	calendar := app.New(storage, log)

	log.Info("calendar is running...")

	server := internalhttp.NewServer(cfg.Server.Host, cfg.Server.Port, calendar, log)
	if err := server.Start(ctx); err != nil {
		log.WithError(err).Fatal("failed to start http server")
	}

	if err := server.Stop(ctx); err != nil {
		log.WithError(err).Error("failed to stop http server")
	}
}
