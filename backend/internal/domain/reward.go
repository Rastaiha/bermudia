package domain

import (
	"math/rand"
)

var rewardSources = map[string]func(player Player) Player{
	"edu1": func(player Player) Player {
		return giveWorthOfCoins(player, 10, false)
	},
	"edu2": func(player Player) Player {
		return giveWorthOfCoins(player, 20, false)
	},
	"edu3": func(player Player) Player {
		return giveWorthOfCoins(player, 30, false)
	},
	"edu4": func(player Player) Player {
		return giveWorthOfCoins(player, 40, false)
	},
	"edu5": func(player Player) Player {
		return giveWorthOfCoins(player, 50, false)
	},
	"edu6": func(player Player) Player {
		return giveWorthOfCoins(player, 60, false)
	},
}

func IsValidRewardSource(rewardSource string) bool {
	if rewardSource == "" {
		return true
	}
	_, ok := rewardSources[rewardSource]
	return ok
}

type rewardType struct {
	kind   string
	value  int32
	weight int
}

var rewardTypes = []rewardType{
	{
		kind:   CostItemTypeCoin,
		value:  1,
		weight: 40,
	},
	{
		kind:   CostItemTypeBlueKey,
		value:  20,
		weight: 6,
	},
	{
		kind:   CostItemTypeRedKey,
		value:  30,
		weight: 3,
	},
	{
		kind:   CostItemTypeGoldenKey,
		value:  50,
		weight: 1,
	},
}

func giveWorthOfCoins(player Player, worthOfCoins int32, allowCoins bool) Player {
	remaining := worthOfCoins

	for remaining > 0 {
		// filter rewards that can still fit in remaining
		var validRewards []rewardType
		for _, r := range rewardTypes {
			if r.value <= remaining {
				validRewards = append(validRewards, r)
			}
		}

		if len(validRewards) == 0 {
			break
		}

		// special case
		if len(validRewards) == 1 && validRewards[0].value == 1 && validRewards[0].kind == CostItemTypeCoin {
			if allowCoins {
				player.Coins += remaining
			}
			break
		}

		// pick with weighted randomness
		totalWeight := 0
		for _, r := range validRewards {
			totalWeight += r.weight
		}

		choice := rand.Intn(totalWeight)
		var pick rewardType
		for _, r := range validRewards {
			if choice < r.weight {
				pick = r
				break
			}
			choice -= r.weight
		}

		switch pick.kind {
		case CostItemTypeCoin:
			if allowCoins {
				player.Coins++
			}
		case CostItemTypeBlueKey:
			player.BlueKeys++
		case CostItemTypeRedKey:
			player.RedKeys++
		case CostItemTypeGoldenKey:
			player.GoldenKeys++
		}
		remaining -= pick.value
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
		return giveWorthOfCoins(player, 20, true), true
	case PoolMedium:
		return giveWorthOfCoins(player, 40, true), true
	case PoolHard:
		return giveWorthOfCoins(player, 80, true), true
	}
	return player, false
}
