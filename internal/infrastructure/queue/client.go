package queue

import (
	"fmt"

	"github.com/anggggggggggg/boilerplate-clean-and-ddd/config"
	"github.com/hibiken/asynq"
)

// redisOpt mengkonversi QueueConfig menjadi asynq.RedisClientOpt.
func redisOpt(cfg *config.QueueConfig) asynq.RedisClientOpt {
	return asynq.RedisClientOpt{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	}
}

// NewAsynqClient membuat Asynq client untuk memproduksi (enqueue) task.
// Client ini digunakan oleh use case melalui queue implementation.
func NewAsynqClient(cfg *config.QueueConfig) *asynq.Client {
	return asynq.NewClient(redisOpt(cfg))
}

// NewAsynqServer membuat Asynq server untuk mengkonsumsi (process) task.
// Server ini dijalankan sebagai goroutine di background saat aplikasi start.
func NewAsynqServer(cfg *config.QueueConfig) *asynq.Server {
	concurrency := cfg.Concurrency
	if concurrency <= 0 {
		concurrency = 10 // default
	}

	return asynq.NewServer(
		redisOpt(cfg),
		asynq.Config{
			Concurrency: concurrency,
		},
	)
}
