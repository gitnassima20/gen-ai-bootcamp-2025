package main

import (
	"adventure-mud-game/config"
	"adventure-mud-game/internal/game"
	"adventure-mud-game/services"
	"fmt"
	"html/template"

	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func main() {
	engine := html.New("./templates", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	fmt.Println("Server running on http://localhost:3000")
	// handlers.RegisterRoutes(app)
	// Load AWS configuration
	cfg := config.LoadAWSConfig()

	// Create Bedrock runtime client
	client := bedrockruntime.NewFromConfig(cfg)

	// Create Bedrock service
	bedrock := services.NewBedrockService(client, "amazon.nova-lite-v1:0")

	// Invoke the Bedrock service
	response, err := bedrock.Invoke("Explain how to play a Mud game")
	if err != nil {
		fmt.Println("Error invoking Bedrock:", err)
	} else {
		fmt.Println("Bedrock response:", response)
	}

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

	// Handle player commands
	app.Post("/command", func(c *fiber.Ctx) error {
		// Get the player's command from the form input
		command := c.FormValue("command")

		// Call the command handler to process the input
		response := game.HandleCommand(command)

		return c.Render("index", fiber.Map{
			"Title":   "The Forgotten Shrine",
			"Message": response,
		})
	})

	app.Listen(":3000")
}
