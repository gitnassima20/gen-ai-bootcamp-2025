package game

import (
	"strings"
)

var CurrentRoom = WorldMap["entrance"]

var directionMap = map[string]string{
	"↑": "north",
	"↓": "south",
	"→": "east",
	"←": "west",
}

func HandleCommand(command string) string {
	// Normalize the command
	command = strings.TrimSpace(command)

	// Check if the command is a movement command
	if direction, ok := directionMap[command]; ok {
		return Move(direction)
	}

	// move case
	if strings.HasPrefix(command, "移動") {
		parts := strings.Fields(command)
		if len(parts) == 2 {
			directionSymbol := parts[1]
			direction, ok := directionMap[directionSymbol]
			if !ok {
				return "Direction not available"
			}
			return Move(direction)
		}
		return "Use arrows to move"
	}

	switch command {
	case "助けて": // "Help"
		return getHelp()
	case "見る": // "Look" in Japanese
		return "You look around. The world is a mixture of light and dark."
	case "持ち物": // "Inventory" in Japanese
		return "Your inventory contains: torch, flower."
	case "話す": // "Talk" in Japanese
		return TalkToNPC()
	default:
		return "Unknow command"
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
- 移動 [direction]: Move in a direction
`
}
