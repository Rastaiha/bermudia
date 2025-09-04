package domain

type IslandHeader struct {
	ID       string `json:"id"`
	BookID   string `json:"bookId"`
	FromPool bool   `json:"fromPool"`
}

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
	Treasures  []Treasure      `json:"treasures"`
}

type BookComponent struct {
	IFrame   *IslandIFrame `json:"iframe,omitempty"`
	Question *Question     `json:"question,omitempty"`
}

type IslandContent struct {
	Components []IslandComponent `json:"components"`
	Treasures  []IslandTreasure  `json:"treasures"`
}

type IslandComponent struct {
	IFrame *IslandIFrame `json:"iframe,omitempty"`
	Input  *IslandInput  `json:"input,omitempty"`
}

type IslandIFrame struct {
	Url string `json:"url"`
}

type IslandTreasure struct {
	ID       string `json:"id"`
	Unlocked bool   `json:"unlocked"`
}

func GetIslandTreasureOfUserTreasure(treasure UserTreasure) IslandTreasure {
	return IslandTreasure{
		ID:       treasure.TreasureID,
		Unlocked: treasure.Unlocked,
	}
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

const (
	PoolEasy   = "easy"
	PoolMedium = "medium"
	PoolHard   = "hard"
)

func IsPoolIdValid(poolId string) bool {
	return poolId == PoolEasy || poolId == PoolMedium || poolId == PoolHard
}

type TerritoryPoolSettings struct {
	Easy   int32 `json:"easy"`
	Medium int32 `json:"medium"`
	Hard   int32 `json:"hard"`
}

func (t TerritoryPoolSettings) TotalCount() int32 {
	return t.Easy + t.Medium + t.Hard
}
