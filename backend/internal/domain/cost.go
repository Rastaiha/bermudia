package domain

type Cost struct {
	Items []CostItem `json:"items"`
}

const (
	CostItemTypeFuel      = "fuel"
	CostItemTypeCoin      = "coin"
	CostItemTypeBlueKey   = "blueKey"
	CostItemTypeRedKey    = "redKey"
	CostItemTypeGoldenKey = "goldenKey"
)

type CostItem struct {
	Type   string `json:"type"`
	Amount int32  `json:"amount"`
}

func canAfford(player Player, cost Cost) bool {
	_, ok := deductCost(player, cost)
	return ok
}

func getItemField(player *Player, itemType string) *int32 {
	switch itemType {
	case CostItemTypeFuel:
		return &player.Fuel
	case CostItemTypeCoin:
		return &player.Coins
	case CostItemTypeBlueKey:
		return &player.BlueKeys
	case CostItemTypeRedKey:
		return &player.RedKeys
	case CostItemTypeGoldenKey:
		return &player.GoldenKeys
	default:
		return nil
	}
}

func deductCost(player Player, cost Cost) (Player, bool) {
	for _, o := range cost.Items {
		field := getItemField(&player, o.Type)
		if field == nil {
			return player, false
		}
		if *field < o.Amount {
			return player, false
		}
		*field -= o.Amount
	}
	return player, true
}

func addCost(player Player, cost Cost) Player {
	for _, o := range cost.Items {
		field := getItemField(&player, o.Type)
		if field != nil {
			*field += o.Amount
		}
	}
	return player
}

func Diff(old, updated Player) Cost {
	result := Cost{}
	for _, item := range []string{
		CostItemTypeFuel,
		CostItemTypeCoin,
		CostItemTypeBlueKey,
		CostItemTypeRedKey,
		CostItemTypeGoldenKey,
	} {
		o := getItemField(&old, item)
		u := getItemField(&updated, item)
		if o != nil && u != nil && *o != *u {
			result.Items = append(result.Items, CostItem{
				Type:   item,
				Amount: *u - *o,
			})
		}
	}
	return result
}
