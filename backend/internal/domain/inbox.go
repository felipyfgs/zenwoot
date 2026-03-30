package domain

import "time"

type ChannelType string

const (
	ChannelWhatsApp ChannelType = "whatsapp"
	ChannelTelegram ChannelType = "telegram"
	ChannelSignal   ChannelType = "signal"
	ChannelEmail    ChannelType = "email"
)

type InboxStatus string

const (
	InboxStatusActive       InboxStatus = "active"
	InboxStatusInactive     InboxStatus = "inactive"
	InboxStatusConnecting   InboxStatus = "connecting"
	InboxStatusDisconnected InboxStatus = "disconnected"
)

type Inbox struct {
	ID          string      `json:"id"`
	AccountID   string      `json:"accountId"`
	Name        string      `json:"name"`
	ChannelType ChannelType `json:"channelType"`
	ChannelID   string      `json:"channelId"`
	Status      InboxStatus `json:"status"`
	Settings    any         `json:"settings"`
	CreatedAt   time.Time   `json:"createdAt"`
	UpdatedAt   time.Time   `json:"updatedAt"`
}
