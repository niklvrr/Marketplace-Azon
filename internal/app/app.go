package app

import (
	"errors"
	"fmt"
	"github.com/niklvrr/myMarketplace/internal/rdb"
	"log"
	"log/slog"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"github.com/niklvrr/myMarketplace/internal/api/router"
	"github.com/niklvrr/myMarketplace/internal/config"
	"github.com/niklvrr/myMarketplace/internal/db"
	"github.com/niklvrr/myMarketplace/pkg/logger"
)

func Run() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	lgr := logger.NewLog(cfg.App.Env)

	dbUrl := cfg.Database.Url

	db.NewDB(dbUrl, lgr)
	defer db.Db.Close()

	mustRunMigrations(dbUrl, lgr)

	rdb.NewRDB(cfg.Cache.Address, lgr)

	r := router.NewRouter(db.Db, cfg.JWT)
	lgr.Info("Starting server")

	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		lgr.Error("Error starting server", err)
	}

	lgr.Info("Server stopped")
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
