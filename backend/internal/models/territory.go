package models

// Island represents a single island in a territory
type Island struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	X         float64 `json:"x"`
	Y         float64 `json:"y"`
	Width     float64 `json:"width"`
	Height    float64 `json:"height"`
	IconAsset string  `json:"iconAsset"`
}

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
	Islands         []Island `json:"islands"`
	Edges           []Edge   `json:"edges"`
}

// APIResponse is the generic response format for all API endpoints
type APIResponse struct {
	OK     bool        `json:"ok"`
	Error  string      `json:"error,omitempty"`
	Result interface{} `json:"result,omitempty"`
}
