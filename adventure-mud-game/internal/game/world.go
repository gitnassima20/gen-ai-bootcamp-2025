package game

type Room struct {
	Name        string
	Description string
	Exits       map[string]string // Directions -> Room Names
	Items       []string
	NPCs        []string
}

var World = map[string]*Room{
	"forest": {
		Name:        "Dark Forest",
		Description: "You are in a dark forest. You hear rustling nearby.",
		Exits: map[string]string{
			"north": "clearing",
			"south": "cave",
		},
		Items: []string{"torch"},
		NPCs:  []string{"旅人"},
	},
	"clearing": {
		Name:        "Sunny Clearing",
		Description: "A beautiful clearing with a small pond.",
		Exits: map[string]string{
			"south": "forest",
		},
		Items: []string{},
		NPCs:  []string{},
	},
}
