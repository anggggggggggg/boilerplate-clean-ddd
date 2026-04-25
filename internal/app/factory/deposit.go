package factory

import (
	deliveryHttp "github.com/anggggggggggg/boilerplate-clean-and-ddd/internal/delivery/http"
	"github.com/anggggggggggg/boilerplate-clean-and-ddd/internal/infrastructure/cache"
	"github.com/anggggggggggg/boilerplate-clean-and-ddd/internal/infrastructure/queue"
	"github.com/anggggggggggg/boilerplate-clean-and-ddd/internal/infrastructure/repository"
	"github.com/anggggggggggg/boilerplate-clean-and-ddd/internal/usecase"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func NewDepositFactory(db *gorm.DB, redisClient *redis.Client, asynqClient *asynq.Client) *deliveryHttp.DepositHandler {
	// ── Layer 1: Infrastructure ──────────────────────────────────────────────
	// Repository: implementasi deposit.Repository menggunakan GORM
	depositRepo := repository.NewDepositRepository(db)

	// Cache: implementasi usecase.DepositCache menggunakan Redis
	depositCache := cache.NewDepositCache(redisClient)

	// Queue: implementasi usecase.DepositQueue menggunakan Asynq
	depositQueue := queue.NewDepositQueue(asynqClient)

	// ── Layer 2: Use Case ────────────────────────────────────────────────────
	// Use case menerima interface, tidak tahu implementasi konkret-nya
	depositUseCase := usecase.NewDepositUseCase(depositRepo, depositCache, depositQueue)

	// ── Layer 3: Delivery (HTTP Handler) ─────────────────────────────────────
	depositHandler := deliveryHttp.NewDepositHandler(depositUseCase)

	return depositHandler
}
