package domain

import (
	"slices"
)

const (
	travelFuelConsumption = 1
	fuelTankCapacity      = 15
)

type Player struct {
	UserId      int32  `json:"-"`
	AtTerritory string `json:"atTerritory"`
	AtIsland    string `json:"atIsland"`
	Fuel        int32  `json:"fuel"`
	FuelCap     int32  `json:"fuelCap"`
}

const (
	PlayerUpdateEventTravel = "travel"
)

type PlayerUpdateEvent struct {
	Reason string  `json:"reason"`
	Player *Player `json:"player"`
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
