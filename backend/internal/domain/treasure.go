package domain

import (
	"time"
)

var (
	ErrTreasureNotFound = Error{
		reason: ErrorReasonResourceNotFound,
		text:   "treasure not found",
	}

	ErrUserTreasureNotFound = Error{
		reason: ErrorReasonResourceNotFound,
		text:   "user treasure not found",
	}
)

type Treasure struct {
	ID     string
	BookID string
}

type UserTreasure struct {
	UserId     int32
	TreasureID string
	Unlocked   bool
	Cost       Cost
	AltCost    Cost
	Reward     *Cost
	UpdatedAt  time.Time
}

type UnlockTreasureCheckResult struct {
	Feasible      bool   `json:"feasible"`
	CanPayCost    bool   `json:"canPayCost"`
	Cost          Cost   `json:"cost"`
	CanPayAltCost bool   `json:"canPayAltCost"`
	AltCost       Cost   `json:"altCost"`
	Reason        string `json:"reason,omitempty"`
}

func UnlockTreasureCheck(player Player, treasure Treasure, userTreasure UserTreasure, currentIslandBook string) (result UnlockTreasureCheckResult) {
	result.Cost = userTreasure.Cost
	result.AltCost = userTreasure.AltCost
	result.CanPayCost = canAfford(player, userTreasure.Cost)
	result.CanPayAltCost = canAfford(player, userTreasure.AltCost)

	if currentIslandBook == "" || treasure.BookID != currentIslandBook {
		result.Reason = "شما باید وارد سیاره شوید تا بتوانید گنج را باز کنید."
		return
	}

	if userTreasure.Unlocked {
		result.Reason = "این گنج قبلاً باز شده است."
		return
	}

	if !result.CanPayCost && !result.CanPayAltCost {
		result.Reason = "شما کلید کافی برای باز کردن این گنج ندارید. "
		return
	}

	result.Feasible = true
	return
}

func UnlockTreasure(player Player, treasure Treasure, userTreasure UserTreasure, currentIslandBook string, chosenCost string) (*PlayerUpdateEvent, UserTreasure, error) {
	check := UnlockTreasureCheck(player, treasure, userTreasure, currentIslandBook)
	if !check.Feasible {
		return nil, UserTreasure{}, Error{
			text:   check.Reason,
			reason: ErrorReasonRuleViolation,
		}
	}

	var couldAfford bool
	switch chosenCost {
	case "":
		player, couldAfford = deductCost(player, userTreasure.Cost)
	case "alt":
		player, couldAfford = deductCost(player, userTreasure.AltCost)
	default:
		return nil, UserTreasure{}, Error{
			reason: ErrorReasonRuleViolation,
			text:   "invalid chosenCost",
		}
	}
	if !couldAfford {
		return nil, UserTreasure{}, Error{
			text:   "شما دارایی کافی برای بازکردن گنج با این روش را ندارید.",
			reason: ErrorReasonRuleViolation,
		}
	}

	reward := getRewardOfTreasure(userTreasure)
	player = addCost(player, reward)
	userTreasure.Unlocked = true
	userTreasure.Reward = &reward
	return &PlayerUpdateEvent{
		Reason: PlayerUpdateEventUnlockTreasure,
		Player: &player,
	}, userTreasure, nil
}
