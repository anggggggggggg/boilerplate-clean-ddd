package routes

import "github.com/gofiber/fiber/v2"

func NewHealtzRoutes(router fiber.Router) {
	router.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})
}
