package domain

import (
	"errors"
	"math"
	"slices"
	"time"
)

const (
	fuelTankCapacity      = 15
	initialFuelAmount     = fuelTankCapacity / 2
	travelFuelConsumption = 1
	initialCoinsAmount    = 100
	refuelCoinCostPerUnit = 10
	anchoringCoinCost     = 20
)

type Player struct {
	UserId      int32     `json:"-"`
	AtTerritory string    `json:"atTerritory"`
	AtIsland    string    `json:"atIsland"`
	Anchored    bool      `json:"anchored"`
	Fuel        int32     `json:"fuel"`
	FuelCap     int32     `json:"fuelCap"`
	Coins       int32     `json:"coins"`
	UpdatedAt   time.Time `json:"-"`
}

type FullPlayer struct {
	Player
	KnowledgeBars []KnowledgeBar `json:"knowledgeBars"`
}

const (
	PlayerUpdateEventInitial    = "initial"
	PlayerUpdateEventTravel     = "travel"
	PlayerUpdateEventRefuel     = "refuel"
	PlayerUpdateEventCorrection = "correction"
	PlayerUpdateEventAnchor     = "anchor"
)

type PlayerUpdateEvent struct {
	Reason string
	Player *Player
}

type FullPlayerUpdateEvent struct {
	Reason string      `json:"reason"`
	Player *FullPlayer `json:"player"`
}

func NewPlayer(userId int32, startingTerritory *Territory) Player {
	return Player{
		UserId:      userId,
		AtTerritory: startingTerritory.ID,
		AtIsland:    startingTerritory.StartIsland,
		Anchored:    true,
		Fuel:        initialFuelAmount,
		FuelCap:     fuelTankCapacity,
		Coins:       initialCoinsAmount,
	}
}

type Cost struct {
	Items []CostItem `json:"items"`
}

const (
	CostItemTypeFuel = "fuel"
	CostItemTypeCoin = "coin"
)

type CostItem struct {
	Type   string `json:"type"`
	Amount int32  `json:"amount"`
}

func canAfford(player Player, cost Cost) bool {
	_, ok := deduceCost(player, cost)
	return ok
}

func deduceCost(player Player, cost Cost) (Player, bool) {
	for _, o := range cost.Items {
		switch o.Type {
		case CostItemTypeFuel:
			if player.Fuel < o.Amount {
				return player, false
			}
			player.Fuel -= o.Amount
		case CostItemTypeCoin:
			if player.Coins < o.Amount {
				return player, false
			}
			player.Coins -= o.Amount
		default:
			return player, false
		}
	}
	return player, true
}

type TravelCheckResult struct {
	Feasible bool `json:"feasible"`
	// Deprecated
	FuelCost   int32  `json:"fuelCost"`
	TravelCost Cost   `json:"travelCost"`
	Reason     string `json:"reason,omitempty"`
}

func TravelCheck(player Player, fromIsland, toIsland string, territory *Territory) (result TravelCheckResult) {
	result = TravelCheckResult{
		Feasible:   false,
		FuelCost:   travelFuelConsumption,
		TravelCost: Cost{Items: []CostItem{{Type: CostItemTypeFuel, Amount: travelFuelConsumption}}},
	}
	if player.AtIsland != fromIsland {
		result.Reason = "شما در جزیره اعلامی قرار ندارید."
		return
	}
	if !slices.ContainsFunc(territory.Edges, func(edge Edge) bool {
		return (edge.From == fromIsland && edge.To == toIsland) ||
			(edge.From == toIsland && edge.To == fromIsland)
	}) {
		result.Reason = "مسیری بین جزیره کنونی و جزیره مقصد وجود ندارد."
		return
	}
	if !canAfford(player, result.TravelCost) {
		result.Reason = "سوخت شما برای سفر کافی نیست."
		return
	}
	result.Feasible = true
	return
}

func Travel(player Player, fromIsland, toIsland string, territory *Territory) (*PlayerUpdateEvent, error) {
	check := TravelCheck(player, fromIsland, toIsland, territory)
	if !check.Feasible {
		return nil, Error{
			reason: ErrorReasonRuleViolation,
			text:   check.Reason,
		}
	}

	player.Fuel -= travelFuelConsumption
	player.AtIsland = toIsland
	player.Anchored = false
	return &PlayerUpdateEvent{
		Reason: PlayerUpdateEventTravel,
		Player: &player,
	}, nil
}

type RefuelCheckResult struct {
	MaxAvailableAmount int32  `json:"maxAvailableAmount"`
	CoinCostPerUnit    int32  `json:"coinCostPerUnit"`
	MaxReason          string `json:"maxReason"`
}

func RefuelCheck(player Player, territory *Territory) (result RefuelCheckResult) {
	idx := slices.IndexFunc(territory.RefuelIslands, func(ri RefuelIsland) bool {
		return ri.ID == player.AtIsland
	})
	if idx < 0 {
		result.MaxReason = "شما در حال حاضر در جزیره سوخت‌گیری قرار ندارید."
		return
	}
	result.CoinCostPerUnit = refuelCoinCostPerUnit
	fuelCapBound := player.FuelCap - player.Fuel
	coinBound := int32(math.MaxInt32)
	if result.CoinCostPerUnit > 0 {
		coinBound = player.Coins / result.CoinCostPerUnit
	}
	if coinBound < fuelCapBound {
		result.MaxAvailableAmount = coinBound
		result.MaxReason = "موجودی سکه شما تنها برای خرید این میزان سوخت کافی است."
	} else {
		result.MaxAvailableAmount = fuelCapBound
		result.MaxReason = "باک شما بیش از این مقدار گنجایش ندارد."
	}
	return
}

func Refuel(player Player, territory *Territory, amount int32) (*PlayerUpdateEvent, error) {
	check := RefuelCheck(player, territory)
	if amount <= 0 {
		return nil, Error{
			text:   "Invalid refuel amount",
			reason: ErrorReasonRuleViolation,
		}
	}
	if check.MaxAvailableAmount < amount {
		return nil, Error{
			text:   check.MaxReason,
			reason: ErrorReasonRuleViolation,
		}
	}
	player.Fuel += amount
	player.Coins -= amount * check.CoinCostPerUnit
	return &PlayerUpdateEvent{
		Reason: PlayerUpdateEventRefuel,
		Player: &player,
	}, nil
}

type AnchorCheckResult struct {
	Feasible      bool   `json:"feasible"`
	AnchoringCost Cost   `json:"anchoringCost"`
	Reason        string `json:"reason,omitempty"`
}

func AnchorCheck(player Player, islandID string) (result AnchorCheckResult) {
	result.AnchoringCost = Cost{Items: []CostItem{{Type: CostItemTypeCoin, Amount: anchoringCoinCost}}}
	if player.AtIsland != islandID {
		result.Reason = "باید به جزیره سفر کنید تا بتوانید در آن لنگر بیندازید."
		return
	}
	if player.Anchored {
		result.Reason = "در حال حاضر در این جزیره لنگر انداخته اید."
		return
	}

	if !canAfford(player, result.AnchoringCost) {
		result.Reason = "دارایی شما برای لنگر انداختن کافی نیست."
		return
	}

	result.Feasible = true
	return
}

func Anchor(player Player, islandID string) (*PlayerUpdateEvent, error) {
	check := AnchorCheck(player, islandID)
	if !check.Feasible {
		return nil, Error{
			reason: ErrorReasonRuleViolation,
			text:   check.Reason,
		}
	}
	var ok bool
	player, ok = deduceCost(player, check.AnchoringCost)
	if !ok {
		return nil, errors.New("logical error in anchor")
	}
	player.Anchored = true
	return &PlayerUpdateEvent{
		Reason: PlayerUpdateEventAnchor,
		Player: &player,
	}, nil
}
