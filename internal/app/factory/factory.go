package factory

import (
	"github.com/anggggggggggg/boilerplate-clean-and-ddd/internal/delivery/http"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Factory adalah pusat dependency injection (DI container) aplikasi.
// Semua dependency di-wire di sini mengikuti urutan dari dalam ke luar:
//
//	Infrastructure → Use Case → Delivery
type Factory struct {
	Deposit *http.DepositHandler
}

// NewFactory menerima koneksi-koneksi dari luar (DB, Redis, Asynq) dan
// merakit semua layer untuk setiap domain.
func NewFactory(db *gorm.DB, redisClient *redis.Client, asynqClient *asynq.Client) *Factory {
	return &Factory{
		Deposit: NewDepositFactory(db, redisClient, asynqClient),
	}
}
