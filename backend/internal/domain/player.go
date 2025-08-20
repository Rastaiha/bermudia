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

func Travel(player Player, fromIsland, toIsland string, territory *Territory) (*PlayerUpdateEvent, error) {
	if player.AtIsland != fromIsland {
		return nil, Error{
			text:   "شما در جزیره اعلامی قرار ندارید.",
			reason: ErrorReasonResourceNotFound,
		}
	}
	if !slices.ContainsFunc(territory.Edges, func(edge Edge) bool {
		return (edge.From == fromIsland && edge.To == toIsland) ||
			(edge.From == toIsland && edge.To == fromIsland)
	}) {
		return nil, Error{
			text:   "مسیری بین این دو جزیره وجود ندارد.",
			reason: ErrorReasonRuleViolation,
		}
	}
	if player.Fuel-travelFuelConsumption < 0 {
		return nil, Error{
			text:   "سوخت شما برای سفر کافی نیست.",
			reason: ErrorReasonRuleViolation,
		}
	}

	player.Fuel -= travelFuelConsumption
	player.AtIsland = toIsland
	return &PlayerUpdateEvent{
		Reason: PlayerUpdateEventTravel,
		Player: &player,
	}, nil
}
