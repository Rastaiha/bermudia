package domain

type RewardSource struct {
	Type      RewardSourceType
	Constant  int32
	MinRandom int32
	MaxRandom int32
}

type RewardSourceType int

const (
	RewardSourceTypeKnowledge RewardSourceType = 0
)

var RewardSources = map[int]RewardSource{
	1: {
		Type:     RewardSourceTypeKnowledge,
		Constant: 100,
	},
}
