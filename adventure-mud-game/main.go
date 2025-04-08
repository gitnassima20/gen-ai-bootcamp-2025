package main

import (
	"adventure-mud-game/handlers"
	"adventure-mud-game/internal/game"
	"fmt"
	"html/template"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func main() {
	engine := html.New("./templates", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	fmt.Println("Server running on http://localhost:3000")
	handlers.RegisterRoutes(app)

	app.Get("/", func(c *fiber.Ctx) error {
		startRoom := game.WorldMap["forest"]

		return c.Render("index", fiber.Map{
			"Title":       "Welcome the Mud Game",
			"Scene":       template.HTML(startRoom.Scene),
			"RoomName":    startRoom.Name,
			"Description": startRoom.Description,
			"Items":       startRoom.Items,
			"NPCs":        startRoom.NPCs,
		})
	})

	app.Listen(":3000")
}
