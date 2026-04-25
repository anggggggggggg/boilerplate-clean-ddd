package routes

import (
	deliveryHttp "github.com/anggggggggggg/boilerplate-clean-and-ddd/internal/delivery/http"
	"github.com/gofiber/fiber/v2"
)

// NewDepositRoutes mendaftarkan semua endpoint deposit ke router.
// Router yang diterima sudah di-wrap middleware (misal JWT) oleh routes.go.
func NewDepositRoutes(router fiber.Router, handler *deliveryHttp.DepositHandler) {
	router.Post("/", handler.CreateDeposit)
	router.Get("/user/:user_id", handler.GetDepositsByUserID)
	router.Get("/:id", handler.GetDeposit)
}
