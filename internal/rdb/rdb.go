package rdb

import (
	"context"
	"log/slog"

	"github.com/redis/go-redis/v9"
)

func NewRDB(addr string, logger *slog.Logger) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		logger.Error("redis db init error", err)
		return nil
	}

	logger.Info("redis db init ok")
	return rdb
}
