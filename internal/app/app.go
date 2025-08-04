package app

import (
	"errors"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"github.com/niklvrr/myMarketplace/internal/config"
	"github.com/niklvrr/myMarketplace/internal/db"
	"github.com/niklvrr/myMarketplace/pkg/logger"
	"log"
	"log/slog"
)

func Run() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	logger := logger.NewLog(cfg.App.Env)

	dbUrl := cfg.Database.Url

	db.NewDB(dbUrl, logger)
	defer db.Db.Close()

	//time.Sleep(20 * time.Second)

	mustRunMigrations(dbUrl, logger)

	// TODO router

}

func mustRunMigrations(dbUrl string, logger *slog.Logger) {
	if dbUrl == "" {
		logger.Error("dbUrl is empty")
		return
	}

	mg, err := migrate.New(
		"file://migrations",
		dbUrl,
	)
	if err != nil {
		logger.Error("migration init err", err)
		return
	}

	if err := mg.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logger.Error("migration run err", err)
		return
	}

	logger.Info("migration run ok")
}
