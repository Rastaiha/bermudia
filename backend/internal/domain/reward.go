package domain

import (
	"math/rand"
	"time"
)

var rewardSources = map[string]func(player Player) Player{
	"edu1": func(player Player) Player {
		return giveRandomWorthOfCoins(player, 10, eduQuestionRewardTypes)
	},
	"edu2": func(player Player) Player {
		return giveRandomWorthOfCoins(player, 20, eduQuestionRewardTypes)
	},
	"edu3": func(player Player) Player {
		return giveRandomWorthOfCoins(player, 30, eduQuestionRewardTypes)
	},
	"edu4": func(player Player) Player {
		return giveRandomWorthOfCoins(player, 40, eduQuestionRewardTypes)
	},
	"edu5": func(player Player) Player {
		return giveRandomWorthOfCoins(player, 50, eduQuestionRewardTypes)
	},
	"edu6": func(player Player) Player {
		return giveRandomWorthOfCoins(player, 60, eduQuestionRewardTypes)
	},
}

func IsValidRewardSource(rewardSource string) bool {
	if rewardSource == "" {
		return true
	}
	_, ok := rewardSources[rewardSource]
	return ok
}

type rewardParam struct {
	worthOfCoins int32
	weight       int
}

var rewardParams = map[string]rewardParam{
	CostItemTypeCoin: {
		worthOfCoins: 1,
		weight:       40,
	},
	CostItemTypeBlueKey: {
		worthOfCoins: 20,
		weight:       6,
	},
	CostItemTypeRedKey: {
		worthOfCoins: 30,
		weight:       3,
	},
	CostItemTypeGoldenKey: {
		worthOfCoins: 50,
		weight:       1,
	},
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

func giveRandomWorthOfCoins(player Player, worthOfCoins int32, rewardTypes []string) Player {
	type rkp struct {
		kind string
		rewardParam
	}

	remaining := worthOfCoins

	for remaining > 0 {
		// filter rewards that can still fit in remaining
		var validRewards []rkp
		for _, t := range rewardTypes {
			p, ok := rewardParams[t]
			if ok && p.worthOfCoins <= remaining {
				validRewards = append(validRewards, rkp{
					kind:        t,
					rewardParam: p,
				})
			}
		}

		if len(validRewards) == 0 {
			break
		}

		// special case
		if len(validRewards) == 1 && validRewards[0].worthOfCoins == 1 {
			if field := getItemField(&player, validRewards[0].kind); field != nil {
				*field += remaining
			}
			break
		}

		// pick with weighted randomness
		totalWeight := 0
		for _, r := range validRewards {
			totalWeight += r.weight
		}

		choice := rand.Intn(totalWeight)
		var pick rkp
		for _, r := range validRewards {
			if choice < r.weight {
				pick = r
				break
			}
			choice -= r.weight
		}

		if field := getItemField(&player, pick.kind); field != nil {
			*field++
		}
		remaining -= pick.worthOfCoins
	}

	return player
}

func GiveRewardOfSource(player Player, rewardSource string) (Player, bool) {
	source, ok := rewardSources[rewardSource]
	if !ok {
		return player, false
	}
	return source(player), true
}

func GiveRewardOfPool(player Player, poolId string) (Player, bool) {
	switch poolId {
	case PoolEasy:
		return giveRandomWorthOfCoins(player, 20, pooledQuestionRewardTypes), true
	case PoolMedium:
		return giveRandomWorthOfCoins(player, 40, pooledQuestionRewardTypes), true
	case PoolHard:
		return giveRandomWorthOfCoins(player, 80, pooledQuestionRewardTypes), true
	}
	return player, false
}

const (
	roughValueOfMasterKey    = 80
	chanceOfGettingMasterKey = 0.3
)

func getRewardOfTreasure(treasure UserTreasure) Cost {
	worthOfCoins := int32(0)
	for _, item := range treasure.Cost.Items {
		worthOfCoins += item.Amount * rewardParams[item.Type].worthOfCoins
	}

	reward := Cost{}
	if worthOfCoins >= roughValueOfMasterKey && rand.Float64() < chanceOfGettingMasterKey {
		reward.Items = append(reward.Items,
			CostItem{
				Type:   CostItemTypeCoin,
				Amount: worthOfCoins - roughValueOfMasterKey,
			},
			CostItem{
				Type:   CostItemTypeMasterKey,
				Amount: 1,
			},
		)
	} else {
		reward.Items = append(reward.Items,
			CostItem{
				Type:   CostItemTypeCoin,
				Amount: worthOfCoins,
			},
		)
	}

	return reward
}

func utc(blue, red, golden int32) Cost {
	cost := Cost{}
	if blue > 0 {
		cost.Items = append(cost.Items, CostItem{Type: CostItemTypeBlueKey, Amount: blue})
	}
	if red > 0 {
		cost.Items = append(cost.Items, CostItem{Type: CostItemTypeRedKey, Amount: red})
	}
	if golden > 0 {
		cost.Items = append(cost.Items, CostItem{Type: CostItemTypeGoldenKey, Amount: golden})
	}
	return cost
}

func GenerateUserTreasure(userId int32, treasureId string) UserTreasure {
	blueC := rand.Float64() * 0.3
	redC := rand.Float64() * (1 - blueC)
	goldenC := rand.Float64() * (1.2 - redC - blueC)

	blue := int32(blueC * 10)
	red := int32(redC * 6)
	golden := int32(goldenC * 4)
	cost := utc(blue, red, golden)

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
