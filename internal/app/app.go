package app

import (
	"log"
	"net"

	"github.com/anggggggggggg/boilerplate-clean-and-ddd/config"
	"github.com/anggggggggggg/boilerplate-clean-and-ddd/internal/app/factory"
	"github.com/anggggggggggg/boilerplate-clean-and-ddd/internal/app/routes"
	"github.com/anggggggggggg/boilerplate-clean-and-ddd/internal/infrastructure/queue"

	"github.com/gofiber/fiber/v2"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// App adalah root struct yang memegang seluruh komponen aplikasi.
type App struct {
	Fiber       *fiber.App
	Config      *config.Config
	asynqServer *asynq.Server
	asynqMux    *asynq.ServeMux
}

// NewApp merakit seluruh aplikasi:
//  1. Buat factory (wiring semua dependency)
//  2. Daftarkan routes ke Fiber
//  3. Daftarkan task handler ke Asynq ServeMux
func NewApp(cfg *config.Config, db *gorm.DB, redisClient *redis.Client, asynqClient *asynq.Client, asynqServer *asynq.Server) *App {
	// Wiring semua dependency via factory
	container := factory.NewFactory(db, redisClient, asynqClient)

	fiberApp := fiber.New()
	routes.NewRoutes(fiberApp, container)

	// Daftarkan semua task handler ke Asynq mux
	// Setiap task type dipetakan ke fungsi processor-nya
	mux := asynq.NewServeMux()
	mux.HandleFunc(queue.TaskDepositCreated, queue.ProcessDepositCreated)

	return &App{
		Fiber:       fiberApp,
		Config:      cfg,
		asynqServer: asynqServer,
		asynqMux:    mux,
	}
}

// Start menjalankan HTTP server dan Asynq worker secara bersamaan.
func (a *App) Start(listener net.Listener) {
	log.Printf("Success start server at port: %s", a.Config.APP.Port)

	// Jalankan Asynq worker di goroutine terpisah
	// Worker akan terus listen task dari Redis queue
	go func() {
		log.Println("Starting Asynq worker...")
		if err := a.asynqServer.Run(a.asynqMux); err != nil {
			log.Printf("Asynq worker error: %s", err.Error())
		}
	}()

	// Jalankan Fiber HTTP server (blocking)
	if err := a.Fiber.Listener(listener); err != nil {
		log.Fatalf("Failed start server with error %s", err.Error())
	}
}

// Shutdown menghentikan Asynq worker dengan graceful shutdown.
func (a *App) Shutdown() {
	a.asynqServer.Shutdown()
}
