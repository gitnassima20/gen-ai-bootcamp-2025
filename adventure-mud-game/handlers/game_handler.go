package handlers

import (
	"text-adventure-mud/internal/game"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/look", func(c *fiber.Ctx) error {
		return c.SendString(game.Look())
	})

	app.Get("/move/:direction", func(c *fiber.Ctx) error {
		direction := c.Params("direction")
		return c.SendString(game.Move(direction))
	})
}
