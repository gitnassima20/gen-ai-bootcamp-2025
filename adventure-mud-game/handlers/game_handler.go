package handlers

import (
	"adventure-mud-game/internal/game"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App) {
	app.Post("/command", func(c *fiber.Ctx) error {
		type CommandRequest struct {
			Command string `json:"command"`
		}

		var req CommandRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid request")
		}

		// Use HandleCommand to process all game commands
		response := game.HandleCommand(req.Command)

		return c.SendString(response)
	})
}
