package game

type Player struct {
	Location  string
	Inventory []string
}

var CurrentPlayer = &Player{
	Location:  "forest",
	Inventory: []string{},
}
