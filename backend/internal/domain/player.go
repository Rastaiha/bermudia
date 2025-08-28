package domain

import (
	"math"
	"slices"
)

const (
	travelFuelConsumption = 1
	fuelTankCapacity      = 15
	refuelCoinCostPerUnit = 0
)

type Player struct {
	UserId      int32  `json:"-"`
	AtTerritory string `json:"atTerritory"`
	AtIsland    string `json:"atIsland"`
	Fuel        int32  `json:"fuel"`
	FuelCap     int32  `json:"fuelCap"`
	Coins       int32  `json:"coins"`
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
		Fuel:        fuelTankCapacity / 2,
		FuelCap:     fuelTankCapacity,
	}
}

type TravelCheckResult struct {
	Feasible bool   `json:"feasible"`
	FuelCost int32  `json:"fuelCost"`
	Reason   string `json:"reason,omitempty"`
}

func TravelCheck(player Player, fromIsland, toIsland string, territory *Territory) (result TravelCheckResult) {
	result = TravelCheckResult{
		Feasible: false,
		FuelCost: travelFuelConsumption,
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
	if player.Fuel-travelFuelConsumption < 0 {
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
