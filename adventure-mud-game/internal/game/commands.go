package game

import (
	"strings"
)

func HandleCommand(command string) string {
	// Normalize the command
	command = strings.ToLower(strings.TrimSpace(command))

	switch command {
	case "助けて": // "Help" in Japanese
		return getHelp()
	case "見る": // "Look" in Japanese
		return "You look around. The world is a mixture of light and dark."
	case "持ち物": // "Inventory" in Japanese
		return "Your inventory contains: torch, flower."
	case "話す": // "Talk" in Japanese
		return TalkToNPC()
	default:
		return "Unknown command. Type '助けて' for a list of commands." // "Unknown command" in English
	}
}

func getHelp() string {
	return `
Available commands:
- 助けて: Show this help message
- 見る: Look around
- 持ち物: Check your inventory
- 話す: Talk to an NPC
- 使用 [item]: Use an item
- 取る [item]: Pick up an item
`
}
