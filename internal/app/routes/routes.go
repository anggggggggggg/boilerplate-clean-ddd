package routes

import (
	"github.com/anggggggggggg/boilerplate-clean-and-ddd/internal/app/factory"
	"github.com/anggggggggggg/boilerplate-clean-and-ddd/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func NewRoutes(app *fiber.App, container *factory.Factory) {
	routerApi := app.Group("/api")

	// Public routes
	healtzRoutes := routerApi.Group("/healtz")
	NewHealtzRoutes(healtzRoutes)

	// Protected routes — semua endpoint di bawah ini wajib membawa JWT
	protectedRoute := routerApi.Group("", middleware.JWTMiddleware())

	// Deposit routes
	depositRouter := protectedRoute.Group("/deposits")
	NewDepositRoutes(depositRouter, container.Deposit)
}
