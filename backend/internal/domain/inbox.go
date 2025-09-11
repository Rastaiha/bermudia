package domain

import (
	"time"
)

type InboxMessageView struct {
	CreatedAt string              `json:"createdAt"`
	Content   InboxMessageContent `json:"content"`
}

type InboxMessage struct {
	ID        string
	UserID    int32
	CreatedAt time.Time
	Content   InboxMessageContent
}

type InboxMessageContent struct {
	NewCorrection    *InboxMessageNewCorrection    `json:"newCorrection,omitempty"`
	OwnOfferAccepted *InboxMessageOwnOfferAccepted `json:"ownOfferAccepted,omitempty"`
}

type InboxEvent struct {
	UserId     int32             `json:"-"`
	NewMessage *InboxMessageView `json:"newMessage,omitempty"`
}

type InboxMessageNewCorrection struct {
	TerritoryID   string          `json:"territoryId"`
	TerritoryName string          `json:"territoryName"`
	IslandID      string          `json:"islandId"`
	IslandName    string          `json:"islandName"`
	InputID       string          `json:"inputId"`
	NewState      SubmissionState `json:"newState"`
	Reward        *Cost           `json:"reward,omitempty"`
}

type InboxMessageOwnOfferAccepted struct {
	Offer TradeOfferView `json:"offer"`
}
