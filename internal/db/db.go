package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
)

var Db *pgxpool.Pool

func NewDB(dbPath string, logger *slog.Logger) {
	if dbPath == "" {
		logger.Error("dbPath is empty")
		return
	}

	var err error
	Db, err = pgxpool.New(context.Background(), dbPath)
	if err != nil {
		logger.Error("db init err", err)
		return
	}

	logger.Info("db init ok")
}
