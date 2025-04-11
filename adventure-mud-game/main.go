package main

import (
	"adventure-mud-game/handlers"
	"adventure-mud-game/internal/game"
	"fmt"
	"html/template"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func main() {
	engine := html.New("./templates", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Serve static files from templates directory
	app.Static("/img", "./templates/img")

	fmt.Println("Server running on http://localhost:3000")
	handlers.RegisterRoutes(app)

	app.Get("/", func(c *fiber.Ctx) error {
		// Set the initial player location
		startRoom := game.WorldMap[game.CurrentPlayer.Location]

		return c.Render("index", fiber.Map{
			"Title":       "Welcome the Mud Game",
			"Scene":       template.HTML(startRoom.Scene),
			"RoomName":    startRoom.Name,
			"Description": startRoom.Description,
			"Items":       startRoom.Items,
			"NPCs":        startRoom.NPCs,
		})
	})

	// Handle player commands
	app.Post("/command", func(c *fiber.Ctx) error {
		// Get the player's command from the form input
		command := c.FormValue("command")

		// Call the command handler to process the input
		response := game.HandleCommand(command)

		// Get the current room after command processing
		currentRoom := game.WorldMap[game.CurrentPlayer.Location]
		fmt.Println("CurrentRoom", currentRoom)

		// Determine the scene name for color scheme
		sceneName := strings.ToLower(strings.ReplaceAll(currentRoom.Name, " ", ""))

		// Render the entire page with updated room details
		return c.Render("index", fiber.Map{
			"Title":       "Into the " + currentRoom.Name,
			"Message":     response,
			"Scene":       template.HTML(currentRoom.Scene),
			"SceneName":   sceneName,
			"RoomName":    currentRoom.Name,
			"Description": currentRoom.Description,
			"Items":       currentRoom.Items,
			"NPCs":        currentRoom.NPCs,
		})
	})

	app.Listen(":3000")
}
