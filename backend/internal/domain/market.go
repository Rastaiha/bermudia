package domain

import (
	"fmt"
	"slices"
	"time"
)

var (
	ErrOfferNotFound = Error{
		reason: ErrorReasonResourceNotFound,
		text:   "offer not found",
	}
)

type TradeOfferView struct {
	ID         string `json:"id"`
	By         string `json:"by"`
	Offered    Cost   `json:"offered"`
	Requested  Cost   `json:"requested"`
	CreatedAt  string `json:"created_at"`
	Acceptable bool   `json:"acceptable"`
}

type TradeOffer struct {
	ID        string    `json:"id"`
	By        int32     `json:"by"`
	Offered   Cost      `json:"offered"`
	Requested Cost      `json:"requested"`
	CreatedAt time.Time `json:"created_at"`
}

func TradeOfferViewForPlayer(player Player, offererUsername string, offer TradeOffer) TradeOfferView {
	return TradeOfferView{
		ID:         offer.ID,
		By:         offererUsername,
		Offered:    offer.Offered,
		Requested:  offer.Requested,
		CreatedAt:  fmt.Sprint(offer.CreatedAt.UnixMilli()),
		Acceptable: isAcceptable(player, offer) == nil,
	}
}

var tradableItems = []string{
	CostItemTypeCoin,
	CostItemTypeBlueKey,
	CostItemTypeRedKey,
	CostItemTypeGoldenKey,
}

type MakeOfferCheckResult struct {
	Feasible      bool   `json:"feasible"`
	TradableItems Cost   `json:"tradableItems"`
	Reason        string `json:"reason,omitempty"`
}

func MakeOfferCheck(player Player, numberOfOpenOffers int) (result MakeOfferCheckResult) {
	for _, i := range tradableItems {
		field := getItemField(&player, i)
		if field != nil {
			maxItem := CostItem{
				Type:   i,
				Amount: *field,
			}
			result.TradableItems.Items = append(result.TradableItems.Items, maxItem)
		}
	}
	if numberOfOpenOffers >= 10 {
		result.Reason = fmt.Sprintf("نمی‌توانید بیش از %d پیشنهاد باز داشته باشید.", numberOfOpenOffers)
		return
	}
	result.Feasible = true
	return
}

func validateAndNormalizeOfferCost(cost Cost) (Cost, error) {
	normalized := make(map[string]CostItem)

	for _, i := range cost.Items {
		if i.Amount == 0 {
			continue
		}
		if i.Amount < 0 {
			return cost, Error{
				reason: ErrorReasonRuleViolation,
				text:   fmt.Sprintf("invalid amount %q", i.Amount),
			}
		}
		if !slices.Contains(tradableItems, i.Type) {
			return cost, Error{
				reason: ErrorReasonRuleViolation,
				text:   fmt.Sprintf("untradable item %q", i.Type),
			}
		}

		current := normalized[i.Type]
		current.Type = i.Type
		current.Amount += i.Amount
		normalized[i.Type] = current
	}
	if len(normalized) == 0 {
		return cost, Error{
			reason: ErrorReasonRuleViolation,
			text:   fmt.Sprintf("empty offer"),
		}
	}

	result := Cost{}
	for _, i := range normalized {
		result.Items = append(result.Items, i)
	}
	slices.SortFunc(result.Items, func(a, b CostItem) int {
		return slices.Index(tradableItems, a.Type) - slices.Index(tradableItems, b.Type)
	})
	return result, nil
}

func MakeOffer(player Player, numberOfOpenOffers int, offered, requested Cost) (*PlayerUpdateEvent, TradeOffer, error) {
	check := MakeOfferCheck(player, numberOfOpenOffers)
	if !check.Feasible {
		return nil, TradeOffer{}, Error{
			reason: ErrorReasonRuleViolation,
			text:   check.Reason,
		}
	}

	var err error
	offered, err = validateAndNormalizeOfferCost(offered)
	if err != nil {
		return nil, TradeOffer{}, err
	}
	requested, err = validateAndNormalizeOfferCost(requested)
	if err != nil {
		return nil, TradeOffer{}, err
	}

	player, ok := deductCost(player, offered)
	if !ok {
		return nil, TradeOffer{}, Error{
			reason: ErrorReasonRuleViolation,
			text:   "دارایی شما برای ثبت این پیشنهاد کافی نیست.",
		}
	}

	tradeOffer := TradeOffer{
		ID:        NewID(ResourceTypeTradeOffer),
		By:        player.UserId,
		Offered:   offered,
		Requested: requested,
		CreatedAt: time.Now().UTC(),
	}

	return &PlayerUpdateEvent{
		Reason: PlayerUpdateEventMakeOffer,
		Player: &player,
	}, tradeOffer, nil
}

func isAcceptable(player Player, offer TradeOffer) error {
	if offer.By == player.UserId {
		return Error{
			reason: ErrorReasonRuleViolation,
			text:   "can't accept your own offer",
		}
	}
	_, ok := deductCost(player, offer.Requested)
	if !ok {
		return Error{
			reason: ErrorReasonRuleViolation,
			text:   "دارایی شما برای قبول این درخواست کافی نیست.",
		}
	}
	return nil
}

func AcceptOffer(acceptor Player, offerer Player, offer TradeOffer) (*PlayerUpdateEvent, *PlayerUpdateEvent, error) {
	err := isAcceptable(acceptor, offer)
	if err != nil {
		return nil, nil, err
	}
	acceptor, ok := deductCost(acceptor, offer.Requested)
	if !ok {
		return nil, nil, Error{
			reason: ErrorReasonRuleViolation,
			text:   "دارایی شما برای قبول این درخواست کافی نیست.",
		}
	}
	offerer = inductCost(offerer, offer.Requested)
	acceptor = inductCost(acceptor, offer.Offered)

	return &PlayerUpdateEvent{
			Reason: PlayerUpdateEventAcceptOffer,
			Player: &acceptor,
		}, &PlayerUpdateEvent{
			Reason: PlayerUpdateEventOwnOfferAccepted,
			Player: &offerer,
		}, nil
}

func DeleteOffer(player Player, offer TradeOffer) (*PlayerUpdateEvent, error) {
	if offer.By != player.UserId {
		return nil, Error{
			reason: ErrorReasonRuleViolation,
			text:   "can't delete an offer that doesn't belong to you",
		}
	}
	player = inductCost(player, offer.Offered)
	return &PlayerUpdateEvent{
		Reason: PlayerUpdateEventOwnOfferDeleted,
		Player: &player,
	}, nil
}
