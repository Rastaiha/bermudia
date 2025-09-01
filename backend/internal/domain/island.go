package domain

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

type Book struct {
	ID         string          `json:"id"`
	Components []BookComponent `json:"components"`
}

type BookComponent struct {
	IFrame   *IslandIFrame `json:"iframe,omitempty"`
	Question *Question     `json:"question,omitempty"`
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
	ID              string          `json:"id"`
	Type            string          `json:"type"`
	Accept          []string        `json:"accept,omitempty"`
	Description     string          `json:"description"`
	SubmissionState SubmissionState `json:"submissionState"`
}

type SubmissionState struct {
	Submittable bool   `json:"submittable"`
	Status      string `json:"status"`
	Filename    string `json:"filename,omitempty"`
	Value       string `json:"value,omitempty"`
	SubmittedAt int64  `json:"submittedAt,omitempty,string"`
}

func PlayerHasAccessToIsland(player Player, islandID string) error {
	if player.AtIsland == islandID && player.Anchored {
		return nil
	}
	return Error{
		reason: ErrorReasonRuleViolation,
		text:   "شما باید در این جزیره لنگر بیندازید تا بتوانید وارد آن شوید.",
	}
}
