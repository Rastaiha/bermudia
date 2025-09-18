package domain

import (
	"errors"
	"fmt"
	"math"
	"time"
)

var (
	ErrInvestSessionNotFound = Error{
		text:   "invest session not found",
		reason: ErrorReasonResourceNotFound,
	}
)

type InvestmentSession struct {
	ID       string
	Text     string
	Resolved bool
	EndAt    time.Time
}

type InvestmentSessionView struct {
	ID    string `json:"id"`
	Text  string `json:"text"`
	EndAt string `json:"endAt"`
}

type InvestmentCheckResult struct {
	Feasible    bool                   `json:"feasible"`
	Reason      string                 `json:"reason,omitempty"`
	Session     *InvestmentSessionView `json:"session,omitempty"`
	Investments []UserInvestment       `json:"investments,omitempty"`
	MaxCoin     int32                  `json:"maxCoin"`
}

type UserInvestment struct {
	SessionID string `json:"-"`
	UserID    int32  `json:"-"`
	Coin      int32  `json:"coin"`
}

func InvestCheck(activeSession *InvestmentSession, investments []UserInvestment, player Player) (result InvestmentCheckResult) {
	result.MaxCoin = player.Coin
	if activeSession == nil {
		result.Reason = "بورس بسته است!"
		return
	}
	result.Session = &InvestmentSessionView{
		ID:    activeSession.ID,
		Text:  activeSession.Text,
		EndAt: fmt.Sprint(activeSession.EndAt.UnixMilli()),
	}
	result.Investments = investments

	if time.Now().UTC().After(activeSession.EndAt.UTC()) {
		result.Reason = "مهلت سرمایه‌گذاری به اتمام رسیده است!"
		return
	}

	if len(investments) > 0 {
		result.Reason = "شما قبلاً یک بار سرمایه‌گذاری کرده اید!"
		return
	}

	if result.MaxCoin <= 0 {
		result.Reason = "شما کلاه کافی برای سرمایه‌گذاری در بورس کلاه‌بهادر ندارید!"
		return
	}

	result.Feasible = true
	return
}

func Invest(session InvestmentSession, investments []UserInvestment, player Player, coinAmount int32) (*PlayerUpdateEvent, UserInvestment, error) {
	check := InvestCheck(&session, investments, player)
	if !check.Feasible {
		return nil, UserInvestment{}, Error{
			text:   check.Reason,
			reason: ErrorReasonRuleViolation,
		}
	}

	if coinAmount <= 0 {
		return nil, UserInvestment{}, Error{
			text:   "coin amount must be positive",
			reason: ErrorReasonRuleViolation,
		}
	}

	player, ok := deductCost(player, Cost{Items: []CostItem{{
		Type:   CostItemTypeCoin,
		Amount: coinAmount,
	}}})

	if !ok {
		return nil, UserInvestment{}, Error{
			text:   "شما کلاه کافی برای این میزان سرمایه‌گذاری ندارید.",
			reason: ErrorReasonRuleViolation,
		}
	}

	ui := UserInvestment{
		SessionID: session.ID,
		UserID:    player.UserId,
		Coin:      coinAmount,
	}

	return &PlayerUpdateEvent{
		Reason: PlayerUpdateEventInvest,
		Player: &player,
	}, ui, nil
}

func ResolveInvestments(session InvestmentSession, investments []UserInvestment, coefficient float64) (map[int32]int32, error) {
	if session.Resolved {
		return nil, errors.New("already resolved")
	}

	diff := session.EndAt.Sub(time.Now().UTC())
	if diff >= 0 {
		return nil, fmt.Errorf("%s is still remaining", diff.String())
	}

	if coefficient < 0 {
		return nil, errors.New("coefficient must not be negative")
	}

	rewards := make(map[int32]int32)

	for _, inv := range investments {
		rewards[inv.UserID] += int32(math.Round(float64(inv.Coin) * coefficient))
	}

	return rewards, nil
}

func GiveInvestmentReward(player Player, coins int32) (PlayerUpdateEvent, bool) {
	player = addCost(player, Cost{Items: []CostItem{{
		Type:   CostItemTypeCoin,
		Amount: coins,
	}}})
	return PlayerUpdateEvent{
		Reason: PlayerUpdateEventInvestReward,
		Player: &player,
	}, coins > 0
}
