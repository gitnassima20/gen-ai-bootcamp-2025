package game

import "fmt"

func Look() string {
	room := World[CurrentPlayer.Location]
	return fmt.Sprintf("%s\n%s\nItems: %v\nNPCs: %v", room.Name, room.Description, room.Items, room.NPCs)
}

func Move(direction string) string {
	room := World[CurrentPlayer.Location]
	if nextRoom, ok := room.Exits[direction]; ok {
		CurrentPlayer.Location = nextRoom
		return "You moved to " + nextRoom
	}
	return "You can't go that way."
}
