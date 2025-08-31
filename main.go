package main

import (
	"embed"
	"fmt"
	"io/fs"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/rs/zerolog"
	"github.com/ziflex/rm-rf-production/internal/api"
	"github.com/ziflex/rm-rf-production/internal/database"
	"github.com/ziflex/rm-rf-production/internal/server"
	"github.com/ziflex/rm-rf-production/pkg/accounts"
	"github.com/ziflex/rm-rf-production/pkg/transactions"
)

//go:embed api/openapi.yaml
var specFile []byte

//go:embed api/ui/*
var uiFS embed.FS

type Config struct {
	Port     int           `env:"PORT" envDefault:"8080"`
	LogLevel zerolog.Level `env:"LOG_LEVEL" envDefault:"trace"`
	DbHost   string        `env:"DB_HOST" envDefault:"localhost"`
	DbPort   int           `env:"DB_PORT" envDefault:"5432"`
	DbName   string        `env:"DB_NAME" envDefault:"mydb"`
	DbUser   string        `env:"DB_USER" envDefault:"user"`
	DbPass   string        `env:"DB_PASS" envDefault:"password"`
}

func main() {
	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("failed to parse env: %+v\n", err)
		os.Exit(1)
	}

	uiSub, err := fs.Sub(uiFS, "api/ui")

	if err != nil {
		fmt.Printf("failed to load embedded ui: %+v\n", err)
		os.Exit(1)
	}

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	db, err := database.New(database.Options{
		Name: cfg.DbName,
		Host: cfg.DbHost,
		Port: cfg.DbPort,
		User: cfg.DbUser,
		Pass: cfg.DbPass,
	})

	if err != nil {
		fmt.Printf("failed to connect to db: %+v\n", err)
		os.Exit(1)
	}

	svr, err := server.NewServer(api.NewHandler(
		accounts.NewService(db, database.NewAccounts()),
		transactions.NewService(db, database.NewTransactions()),
	), server.Options{
		Logger: logger,
		Spec:   specFile,
		UI:     uiSub,
	})

	if err != nil {
		fmt.Printf("failed to create server: %+v\n", err)
		os.Exit(1)
	}

	if err = svr.Run(cfg.Port); err != nil {
		fmt.Printf("server error: %+v\n", err)
		os.Exit(1)
	}
}
