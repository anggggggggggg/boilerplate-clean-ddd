package usecase

import (
	"context"
	"log"

	"github.com/anggggggggggg/boilerplate-clean-and-ddd/internal/domain/deposit"
	"github.com/anggggggggggg/boilerplate-clean-and-ddd/internal/dto"
)

// ---------------------------------------------------------------------------
// Port interfaces (didefinisikan di use case, diimplementasikan di infrastruktur)
// ---------------------------------------------------------------------------

// DepositCache adalah kontrak cache untuk deposit.
// Implementasinya ada di internal/infrastructure/cache — use case tidak tahu Redis.
type DepositCache interface {
	Get(ctx context.Context, id string) (*dto.DepositResponse, error)
	Set(ctx context.Context, d *dto.DepositResponse) error
	Delete(ctx context.Context, id string) error
}

// DepositQueue adalah kontrak antrian untuk deposit.
// Implementasinya ada di internal/infrastructure/queue — use case tidak tahu Asynq.
type DepositQueue interface {
	EnqueueDepositCreated(ctx context.Context, depositID string) error
}

// ---------------------------------------------------------------------------
// Use Case interface
// ---------------------------------------------------------------------------

// DepositUseCase mendefinisikan kontrak application service untuk domain deposit.
type DepositUseCase interface {
	CreateDeposit(ctx context.Context, req *dto.CreateDepositRequest) (*dto.DepositResponse, error)
	GetDeposit(ctx context.Context, id string) (*dto.DepositResponse, error)
	GetDepositsByUserID(ctx context.Context, userID string) ([]*dto.DepositResponse, error)
}

// ---------------------------------------------------------------------------
// Implementation
// ---------------------------------------------------------------------------

type depositUseCase struct {
	repo  deposit.Repository
	cache DepositCache
	queue DepositQueue
}

// NewDepositUseCase menerima semua dependency via injection.
// Urutan parameter: repo (DB) → cache (Redis) → queue (Asynq).
func NewDepositUseCase(repo deposit.Repository, cache DepositCache, queue DepositQueue) DepositUseCase {
	return &depositUseCase{
		repo:  repo,
		cache: cache,
		queue: queue,
	}
}

func (uc *depositUseCase) CreateDeposit(ctx context.Context, req *dto.CreateDepositRequest) (*dto.DepositResponse, error) {
	// Domain layer yang validasi business rule
	d, err := deposit.NewDeposit(req.UserID, req.Amount, req.Note)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Create(ctx, d); err != nil {
		return nil, err
	}

	res := mapToResponse(d)

	// Enqueue job secara asynchronous — jika queue gagal, jangan blok response
	if err := uc.queue.EnqueueDepositCreated(ctx, d.ID); err != nil {
		log.Printf("[UseCase] gagal enqueue deposit_created untuk id=%s: %v", d.ID, err)
	}

	return res, nil
}

func (uc *depositUseCase) GetDeposit(ctx context.Context, id string) (*dto.DepositResponse, error) {
	// Cache hit — langsung return tanpa ke DB
	if cached, err := uc.cache.Get(ctx, id); err == nil {
		log.Printf("[UseCase] cache HIT untuk deposit id=%s", id)
		return cached, nil
	}

	// Cache miss — ambil dari DB
	log.Printf("[UseCase] cache MISS untuk deposit id=%s, fetch dari DB", id)
	d, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	res := mapToResponse(d)

	// Simpan ke cache, error tidak kritis
	if err := uc.cache.Set(ctx, res); err != nil {
		log.Printf("[UseCase] gagal set cache untuk deposit id=%s: %v", id, err)
	}

	return res, nil
}

// GetDepositsByUserID mengambil semua deposit milik user dari DB langsung.
func (uc *depositUseCase) GetDepositsByUserID(ctx context.Context, userID string) ([]*dto.DepositResponse, error) {
	deposits, err := uc.repo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.DepositResponse, 0, len(deposits))
	for _, d := range deposits {
		responses = append(responses, mapToResponse(d))
	}

	return responses, nil
}

// mapToResponse mengkonversi domain entity ke DTO response.
// Mapping ini ada di use case agar domain tidak tahu format presentasi.
func mapToResponse(d *deposit.Deposit) *dto.DepositResponse {
	return &dto.DepositResponse{
		ID:        d.ID,
		UserID:    d.UserID,
		Amount:    d.Amount,
		Status:    string(d.Status),
		Note:      d.Note,
		CreatedAt: d.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: d.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
