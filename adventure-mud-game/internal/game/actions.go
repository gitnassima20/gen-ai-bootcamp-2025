package game

import (
	"fmt"
	"log"

	"adventure-mud-game/config"
	"adventure-mud-game/services"

	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

func TalkToNPC() string {

	// Check if there are NPCs in the current room
	currentRoom, exists := WorldMap[CurrentPlayer.Location]
	if !exists {
		log.Printf("TalkToNPC: Room not found - %s", CurrentPlayer.Location)
		return "You seem to be in an undefined location."
	}

	if len(currentRoom.NPCs) == 0 {
		log.Println("TalkToNPC: No NPCs in the current room")
		return "There are no NPCs to talk to in this area."
	}

	// Select the first NPC
	npc := currentRoom.NPCs[0]
	log.Printf("TalkToNPC: Interacting with NPC - %s", npc)

	// Prepare context-aware prompt
	prompt := fmt.Sprintf(`You are a %s in a room called %s. 
Room Description: %s
Room Items: %v
Room Scene: %s

Generate a short, immersive dialogue that fits this context. 
Be creative and mysterious. 
Speak in a way that is intriguing and adds depth to the game world. 
Provide a hint or a cryptic message that might be useful to the player.`,
		npc,
		currentRoom.Name,
		currentRoom.Description,
		currentRoom.Items,
		currentRoom.Scene)

	log.Printf("TalkToNPC: Generated prompt - %s", prompt)

	// Attempt to generate dialogue with fallback
	response, err := generateDialogue(prompt)
	if err != nil {
		log.Printf("TalkToNPC: Dialogue generation failed - %v", err)
		return fmt.Sprintf("The %s seems lost in thought and doesn't respond.", npc)
	}

	log.Printf("TalkToNPC: Received response - %s", response)
	return response
}

func generateDialogue(prompt string) (string, error) {
	// Load AWS configuration
	cfg, err := config.LoadAWSConfig()
	if err != nil {
		return "", fmt.Errorf("failed to load AWS config: %v", err)
	}

	// Create Bedrock runtime client
	client := bedrockruntime.NewFromConfig(cfg)

	// Create Bedrock service
	bedrockService := services.NewBedrockService(client, "amazon.nova-lite-v1:0")

	// Invoke Bedrock service to generate dialogue
	response, err := bedrockService.Invoke(prompt)
	if err != nil {
		return "", fmt.Errorf("failed to generate dialogue: %v", err)
	}

	return response, nil
}
