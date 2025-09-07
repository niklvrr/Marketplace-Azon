package rdb

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log/slog"
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

	return rdb
}
