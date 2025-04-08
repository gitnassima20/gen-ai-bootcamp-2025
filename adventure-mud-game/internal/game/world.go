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
	"forest": {
		Name:        "Dark Forest",
		Description: "You are in a dark forest. You hear rustling nearby.",
		Exits: map[string]string{
			"north": "clearing",
			"south": "cave",
		},
		Items: []string{"torch", "stone"},
		NPCs:  []string{"traveler", "black cat"},
		Scene: `
<pre class="scene">
🌲🌲🌲🌲🌲
🌲  🧍  🌲
🌲 🐾   🐱 🌲
🌲    🔦    🌲
🌲🌲🌲🌲🌲
</pre>`,
	},
	"clearing": {
		Name:        "Sunny Clearing",
		Description: "A beautiful clearing with a small pond.",
		Exits: map[string]string{
			"south": "forest",
		},
		Items: []string{"flower", "water bottle"},
		NPCs:  []string{"girl"},
		Scene: `
<pre class="scene">
🌿🌸🌞🌸🌿
🌸     🧍    🌸
💧   🌼   💧
🌿🌸🌿🌸🌿
</pre>`,
	},
	"cave": {
		Name:        "Echoing Cave",
		Description: "A cold, damp cave with echoes of dripping water. It's dimly lit.",
		Exits: map[string]string{
			"north": "forest",
		},
		Items: []string{"old scroll", "key"},
		NPCs:  []string{"old man"},
		Scene: `
<pre class="scene">
🪨🪨🪨🪨🪨
🪨   🧓   🪨
💧   🧻   🔑
🪨     🕯️    🪨
🪨🪨🪨🪨🪨
</pre>`,
	},
}
