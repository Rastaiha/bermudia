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

func NewPlayer(userId int32, startingTerritory *Territory) Player {
	return Player{
		UserId:      userId,
		AtTerritory: startingTerritory.ID,
		AtIsland:    startingTerritory.StartIsland,
		Fuel:        fuelTankCapacity / 2,
		FuelCap:     fuelTankCapacity,
	}
}

func Travel(player *Player, fromIsland, toIsland string, territory *Territory) error {
	if player.AtIsland != fromIsland {
		return Error{
			text:   "شما در جزیره اعلامی قرار ندارید.",
			reason: ErrorReasonRuleViolation,
		}
	}
	if !slices.ContainsFunc(territory.Edges, func(edge Edge) bool {
		return edge.From == fromIsland && edge.To == toIsland
	}) {
		return Error{
			text:   "مسیری بین این دو جزیره وجود ندارد.",
			reason: ErrorReasonRuleViolation,
		}
	}
	if player.Fuel-travelFuelConsumption < 0 {
		return Error{
			text:   "سوخت شما برای سفر کافی نیست.",
			reason: ErrorReasonRuleViolation,
		}
	}
	player.Fuel -= travelFuelConsumption
	player.AtIsland = toIsland
	return nil
}
