package database

// Facade tipis di atas service-utils/cache (DRY).

import (
	"context"
	"log/slog"
	"time"

	"github.com/vikikurnia87/service-order/configs"

	sucache "github.com/vikikurnia87/service-utils/cache"
)

var client *sucache.Cache

func InitRedisDatabase(ctx context.Context, logger *slog.Logger) {
	c, err := sucache.New(ctx, sucache.Config{
		Host:            configs.RedisHost,
		Port:            configs.RedisPort,
		Password:        configs.RedisPassword,
		DB:              0,
		PoolSize:        configs.RedisPoolSize,
		MinIdleConns:    configs.RedisMinIdleConns,
		MaxRetries:      configs.RedisMaxRetries,
		MinRetryBackoff: configs.RedisMinRetryBackoffS,
		MaxRetryBackoff: configs.RedisMaxRetryBackoffS,
		DialTimeout:     configs.RedisDialTimeoutS,
		ReadTimeout:     configs.RedisReadTimeoutS,
		WriteTimeout:    configs.RedisWriteTimeoutS,
		PoolTimeout:     configs.RedisPoolTimeoutS,
		Expiration:      configs.RedisExpirationS,
		Env:             configs.ServerEnv,
		Namespace:       configs.RedisNamespace,
	}, logger)
	if err != nil {
		logger.ErrorContext(ctx, "Redis init failed, running without cache", slog.String("error", err.Error()))
		return
	}
	client = c
}

func RedisClose(ctx context.Context, logger *slog.Logger) {
	if client != nil {
		client.Close(ctx, logger)
	}
}

// GetCache mengembalikan instance *cache.Cache untuk dependency injection.
func GetCache() *sucache.Cache {
	return client
}

func RedisSet(ctx context.Context, key string, value any, ttlOptional ...time.Duration) error {
	if client == nil {
		return nil
	}
	return client.Set(ctx, key, value, ttlOptional...)
}

func RedisGet(ctx context.Context, key string, dest any) (bool, error) {
	if client == nil {
		return false, nil
	}
	return client.Get(ctx, key, dest)
}

func Delete(ctx context.Context, key string) error {
	if client == nil {
		return nil
	}
	return client.Delete(ctx, key)
}

func DeleteByPattern(ctx context.Context, pattern string) (int64, error) {
	if client == nil {
		return 0, nil
	}
	return client.DeleteByPattern(ctx, pattern)
}
