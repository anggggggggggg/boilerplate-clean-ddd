package main

import (
	"context"
	"log"

	"github.com/anggggggggggg/boilerplate-clean-and-ddd/config"
	"github.com/anggggggggggg/boilerplate-clean-and-ddd/internal/app"
	"github.com/anggggggggggg/boilerplate-clean-and-ddd/internal/infrastructure/cache"
	"github.com/anggggggggggg/boilerplate-clean-and-ddd/internal/infrastructure/db"
	"github.com/anggggggggggg/boilerplate-clean-and-ddd/internal/infrastructure/queue"
	"github.com/anggggggggggg/boilerplate-clean-and-ddd/pkg/graceful"

	"github.com/jpillora/overseer"
	"github.com/jpillora/overseer/fetcher"
)

// cfg dideklarasikan di package scope agar dapat diakses oleh fungsi program().
// Overseer memanggil program() di child process — cfg harus sudah ter-load sebelumnya.
var cfg *config.Config

func main() {
	var err error
	cfg, err = config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	debug := cfg.APP.Env == "development"

	// Overseer mengelola hot-reload binary tanpa downtime.
	// Listener (TCP port) dikelola oleh overseer dan diteruskan ke program().
	overseer.Run(overseer.Config{
		Program:       program,
		Address:       ":" + cfg.APP.Port,
		Fetcher:       &fetcher.File{Path: cfg.APP.ProgramFile, Interval: 5},
		Debug:         debug,
		RestartSignal: graceful.RestartSignal,
	})
}

// program adalah entry point yang dijalankan oleh overseer di child process.
// Urutan inisialisasi: DB → Redis (cache) → Asynq (queue) → App → Start
func program(state overseer.State) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup graceful shutdown — akan memanggil cancel() saat menerima sinyal OS
	graceful.SetupGracefulShutdown(cancel)

	// ── 1. Database ──────────────────────────────────────────────────────────
	dbConn, err := db.ConnectDB(cfg.DB)
	if err != nil {
		log.Fatalf("Failed to connect database: %s", err.Error())
	}

	// ── 2. Cache (Redis) ─────────────────────────────────────────────────────
	redisClient, err := cache.NewRedisClient(cfg.Redis)
	if err != nil {
		log.Fatalf("Failed to connect Redis: %s", err.Error())
	}

	// ── 3. Queue (Asynq via Redis) ───────────────────────────────────────────
	// Client untuk produce (enqueue) task
	asynqClient := queue.NewAsynqClient(cfg.Queue)
	// Server untuk consume (process) task — dijalankan sebagai worker goroutine
	asynqServer := queue.NewAsynqServer(cfg.Queue)

	// ── 4. Rakit dan jalankan aplikasi ───────────────────────────────────────
	application := app.NewApp(cfg, dbConn, redisClient, asynqClient, asynqServer)
	application.Start(state.Listener)

	// Tunggu sampai context di-cancel (sinyal shutdown)
	<-ctx.Done()

	// ── 5. Graceful shutdown — urutan kebalikan dari inisialisasi ────────────
	application.Shutdown() // hentikan Asynq worker
	_ = asynqClient.Close()
	_ = redisClient.Close()
	db.CloseDB()

	log.Println("Shutting down gracefully...")
	log.Println("Cleanup done. Exiting.")
}
