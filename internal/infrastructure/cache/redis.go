package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/anggggggggggg/boilerplate-clean-and-ddd/config"
	"github.com/redis/go-redis/v9"
)

// NewRedisClient membuat koneksi ke Redis menggunakan konfigurasi dari config.RedisConfig.
// Koneksi di-ping untuk memastikan Redis dapat dijangkau sebelum aplikasi berjalan.
func NewRedisClient(cfg *config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:       fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password:   cfg.Password,
		DB:         cfg.DB,
		MaxRetries: cfg.MaxRetries,
		PoolSize:   cfg.PoolSize,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("gagal konek ke Redis: %w", err)
	}

	return client, nil
}
