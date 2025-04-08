package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"adventure-mud-game/handlers"
)

func main() {
	app := fiber.New()

	fmt.Println("Server running on http://localhost:3000")
	handlers.RegisterRoutes(app)

	app.Listen(":3000")
}
