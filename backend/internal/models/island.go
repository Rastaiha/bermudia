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

type IslandContent struct {
	Components []IslandComponent `json:"components"`
}

type IslandComponent struct {
	IFrame *IslandIFrame `json:"iframe,omitempty"`
	Input  *IslandInput  `json:"input,omitempty"`
}

type IslandIFrame struct {
	Url string `json:"url"`
}

type IslandInput struct {
	ID          string   `json:"id"`
	Type        string   `json:"type"`
	Accept      []string `json:"accept,omitempty"`
	Description string   `json:"description"`
}
