package domain

// Edge represents a connection between two islands
type Edge struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// Territory represents a complete territory with islands and their connections
type Territory struct {
	ID              string         `json:"id"`
	Name            string         `json:"name"`
	BackgroundAsset string         `json:"backgroundAsset"`
	StartIsland     string         `json:"startIsland"`
	Islands         []Island       `json:"islands"`
	Edges           []Edge         `json:"edges"`
	RefuelIslands   []RefuelIsland `json:"refuelIslands"`
}

type RefuelIsland struct {
	ID string `json:"id"`
}
