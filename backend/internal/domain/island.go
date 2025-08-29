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

type IslandInputContent struct {
	Components []*IslandInputComponent `json:"components"`
}

type IslandInputComponent struct {
	ID       string               `json:"id,omitempty"`
	IFrame   *IslandIFrame        `json:"iframe,omitempty"`
	Question *IslandInputQuestion `json:"question,omitempty"`
}

type IslandInputQuestion struct {
	Question
	KnowledgeAmount int32 `json:"knowledgeAmount"`
}

type IslandRawContent struct {
	Components []IslandRawComponent `json:"components"`
}

type IslandRawComponent struct {
	ID       string             `json:"id,omitempty"`
	IFrame   *IslandIFrame      `json:"iframe,omitempty"`
	Question *QuestionComponent `json:"question,omitempty"`
}

type QuestionComponent struct {
	QuestionID string `json:"questionId"`
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

type UserComponent struct {
	IslandID    string
	UserID      int32
	ComponentID string
	ResourceID  string
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
