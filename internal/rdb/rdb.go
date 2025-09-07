package rdb

import (
	"context"
	"log/slog"

	"github.com/redis/go-redis/v9"
)

var CacheDB *redis.Client

func NewRDB(addr string, logger *slog.Logger) {
	CacheDB = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	if err := CacheDB.Ping(context.Background()).Err(); err != nil {
		logger.Error("redis db init error", err)
		return
	}

	logger.Info("redis db init ok")
	return
}
