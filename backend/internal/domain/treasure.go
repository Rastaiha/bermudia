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
	UpdatedAt  time.Time
}

type UnlockTreasureCheckResult struct {
	Feasible bool   `json:"feasible"`
	Cost     Cost   `json:"cost"`
	Reason   string `json:"reason,omitempty"`
}

func UnlockTreasureCheck(player Player, treasure UserTreasure) (result UnlockTreasureCheckResult) {
	result.Cost = treasure.Cost
	if treasure.Unlocked {
		result.Reason = "این گنج قبلاً باز شده است."
		return
	}

	if _, ok := deduceCost(player, treasure.Cost); !ok {
		result.Reason = "شما کلید کافی برای باز کردن این گنج ندارید. "
		return
	}

	result.Feasible = true
	return
}

func UnlockTreasure(player Player, treasure UserTreasure) (*PlayerUpdateEvent, UserTreasure, error) {
	check := UnlockTreasureCheck(player, treasure)
	if !check.Feasible {
		return nil, UserTreasure{}, Error{
			text:   check.Reason,
			reason: ErrorReasonRuleViolation,
		}
	}
	player, ok := deduceCost(player, check.Cost)
	if !ok {
		return nil, UserTreasure{}, Error{
			text:   "شما دارایی کافی ندارید.",
			reason: ErrorReasonRuleViolation,
		}
	}
	player = giveRewardOfTreasure(player, treasure)
	treasure.Unlocked = true
	return &PlayerUpdateEvent{
		Reason: PlayerUpdateEventUnlockTreasure,
		Player: &player,
	}, treasure, nil
}
