package domain

import (
	"fmt"
	"time"
)

var (
	ErrInvestSessionNotFound = Error{
		text:   "invest session not found",
		reason: ErrorReasonResourceNotFound,
	}
)

type InvestmentSession struct {
	ID               string
	Text             string
	Resolved         bool
	RewardThresholds map[int32]float64
	EndAt            time.Time
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
