package domain

import (
	"math/rand"
	"time"
)

var rewardSources = map[string]func(player Player) Player{
	"edu1": func(player Player) Player {
		return giveRandomWorthOfCoins(player, 35, eduQuestionRewardTypes)
	},
	"edu2": func(player Player) Player {
		return giveRandomWorthOfCoins(player, 50, eduQuestionRewardTypes)
	},
	"edu3": func(player Player) Player {
		return giveRandomWorthOfCoins(player, 65, eduQuestionRewardTypes)
	},
	"edu4": func(player Player) Player {
		return giveRandomWorthOfCoins(player, 80, eduQuestionRewardTypes)
	},
	"edu5": func(player Player) Player {
		return giveRandomWorthOfCoins(player, 95, eduQuestionRewardTypes)
	},
	"edu6": func(player Player) Player {
		return giveRandomWorthOfCoins(player, 110, eduQuestionRewardTypes)
	},
}

func IsValidRewardSource(rewardSource string) bool {
	if rewardSource == "" {
		return true
	}
	_, ok := rewardSources[rewardSource]
	return ok
}

var rewardParams = map[string]int32{
	CostItemTypeCoin:      1,
	CostItemTypeBlueKey:   20,
	CostItemTypeRedKey:    30,
	CostItemTypeGoldenKey: 50,
}

var eduQuestionRewardTypes = []string{
	CostItemTypeBlueKey,
	CostItemTypeRedKey,
	CostItemTypeGoldenKey,
}

var pooledQuestionRewardTypes = []string{
	CostItemTypeCoin,
	CostItemTypeBlueKey,
	CostItemTypeRedKey,
	CostItemTypeGoldenKey,
}

var treasureCostRewardTypes = []string{
	CostItemTypeBlueKey,
	CostItemTypeRedKey,
	CostItemTypeGoldenKey,
}

func giveRandomWorthOfCoins(player Player, worthOfCoins int32, rewardTypes []string) Player {
	type rkp struct {
		kind         string
		worthOfCoins int32
	}

	// Build validRewards once before the loop
	var validRewards []rkp
	for _, t := range rewardTypes {
		worthOfCoins, ok := rewardParams[t]
		if ok {
			validRewards = append(validRewards, rkp{
				kind:         t,
				worthOfCoins: worthOfCoins,
			})
		}
	}

	if len(validRewards) == 0 {
		return player
	}

	remaining := worthOfCoins

	for remaining > 0 {
		// Calculate total weight (inverse of worthOfCoins, scaled by 1000)
		totalWeight := 0
		for _, r := range validRewards {
			totalWeight += int(1000 / r.worthOfCoins)
		}

		choice := rand.Intn(totalWeight)
		var pick rkp
		for _, r := range validRewards {
			weight := int(1000 / r.worthOfCoins)
			if choice < weight {
				pick = r
				break
			}
			choice -= weight
		}

		player = addCost(player, Cost{Items: []CostItem{{Type: pick.kind, Amount: 1}}})
		remaining -= pick.worthOfCoins
	}

	return player
}

func giveRewardOfSource(player Player, rewardSource string) (Player, bool) {
	source, ok := rewardSources[rewardSource]
	if !ok {
		return player, false
	}
	return source(player), true
}

func giveRewardOfPool(player Player, poolId string) (Player, bool) {
	switch poolId {
	case PoolEasy:
		return giveRandomWorthOfCoins(player, 50, pooledQuestionRewardTypes), true
	case PoolMedium:
		return giveRandomWorthOfCoins(player, 75, pooledQuestionRewardTypes), true
	case PoolHard:
		return giveRandomWorthOfCoins(player, 125, pooledQuestionRewardTypes), true
	}
	return player, false
}

func GetRewardOfCorrection(player Player, question BookQuestion, correction Correction, pool string, hasPool bool) (*PlayerUpdateEvent, *Cost, bool) {
	if correction.NewStatus != AnswerStatusCorrect && correction.NewStatus != AnswerStatusHalfCorrect {
		return nil, nil, false
	}
	newPlayer, rewarded := giveRewardOfSource(player, question.RewardSource)
	if hasPool {
		var rewarded2 bool
		newPlayer, rewarded2 = giveRewardOfPool(newPlayer, pool)
		rewarded = rewarded || rewarded2
	}
	if !rewarded {
		return nil, nil, false
	}
	reward := Diff(player, newPlayer)
	return &PlayerUpdateEvent{
		Reason: PlayerUpdateEventCorrection,
		Player: &newPlayer,
	}, &reward, true
}

const (
	chanceOfGettingMasterKey = 0.1
	treasureMinCost          = 20
	treasureMaxCost          = 100
)

func getRewardOfTreasure(treasure UserTreasure) Cost {
	worthOfCoins := int32(0)
	for _, item := range treasure.Cost.Items {
		worthOfCoins += item.Amount * rewardParams[item.Type]
	}

	reward := Cost{}
	masterKeyRoughValue := int32(treasureMinCost) + int32(0.6*float64(treasureMaxCost-treasureMinCost))
	if worthOfCoins >= masterKeyRoughValue && rand.Float64() < chanceOfGettingMasterKey {
		worthOfCoins -= masterKeyRoughValue / 2
		reward.Items = append(reward.Items,
			CostItem{
				Type:   CostItemTypeMasterKey,
				Amount: 1,
			},
		)
	}
	if worthOfCoins > 0 {
		reward.Items = append(reward.Items,
			CostItem{
				Type:   CostItemTypeCoin,
				Amount: worthOfCoins,
			},
		)
	}

	return reward
}

func generateTreasureCost(worthOfCoins int32) Cost {
	// Create a dummy player to use with giveRandomWorthOfCoins
	dummyPlayer := Player{}
	updatedPlayer := giveRandomWorthOfCoins(dummyPlayer, worthOfCoins, treasureCostRewardTypes)
	return Diff(Player{}, updatedPlayer)
}

func GenerateUserTreasure(userId int32, treasureId string) UserTreasure {
	costValue := int32(treasureMinCost + rand.Intn(treasureMaxCost-treasureMinCost+1))
	cost := generateTreasureCost(costValue)

	reward := &Cost{}
	if len(cost.Items) > 0 {
		reward = nil
	}
	return UserTreasure{
		UserId:     userId,
		TreasureID: treasureId,
		Unlocked:   len(cost.Items) == 0,
		Cost:       cost,
		AltCost: Cost{
			Items: []CostItem{
				{
					Type:   CostItemTypeMasterKey,
					Amount: 1,
				},
			},
		},
		Reward:    reward,
		UpdatedAt: time.Now().UTC(),
	}
}
