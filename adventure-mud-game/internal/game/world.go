package game

type Room struct {
	Name        string
	Description string
	Exits       map[string]string // Directions -> Room Names
	Items       []string
	Scene       string
	NPCs        []string
}

var WorldMap = map[string]*Room{
	"entrance": {
		Name:        "Entrance",
		Description: "You are at the entrance of the Isekai World",
		Exits: map[string]string{
			"north": "forest",
		},
		Items: []string{"door"},
		NPCs:  []string{},
		Scene: `<img src="/templates/img/entrance.jpg" alt="Entrance" class="scene-image">`,
	},
	"forest": {
		Name:        "Dark Forest",
		Description: "You are in a dark forest. You hear rustling nearby.",
		Exits: map[string]string{
			"north": "clearing",
			"south": "cave",
		},
		Items: []string{"torch"},
		NPCs:  []string{"traveler", "black cat"},
		Scene: `<img src="/templates/img/forest.jpg" alt="Dark Forest" class="scene-image">`,
	},
	"clearing": {
		Name:        "Sunny Clearing",
		Description: "A beautiful clearing with a small pond.",
		Exits: map[string]string{
			"south": "forest",
		},
		Items: []string{"flower", "water bottle"},
		NPCs:  []string{"wolf"},
		Scene: `<img src="/templates/img/clearing.jpg" alt="Sunny Clearing" class="scene-image">`,
	},
	"cave": {
		Name:        "Echoing Cave",
		Description: "A cold, damp cave with echoes of dripping water. It's dimly lit.",
		Exits: map[string]string{
			"north": "forest",
		},
		Items: []string{"sword", "key"},
		NPCs:  []string{"old woman"},
		Scene: `<img src="/templates/img/cave.jpg" alt="Echoing Cave" class="scene-image">`,
	},
}
