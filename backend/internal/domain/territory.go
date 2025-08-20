package domain

// Edge represents a connection between two islands
type Edge struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// Territory represents a complete territory with islands and their connections
type Territory struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	BackgroundAsset string   `json:"backgroundAsset"`
	FirstIsland     *string  `json:"firstIsland,omitempty"`
	Islands         []Island `json:"islands"`
	Edges           []Edge   `json:"edges"`
}
