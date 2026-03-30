package whatsapp

import (
	"encoding/json"
	"time"
)

// wzap API response wrapper
type wzapAPIResponse struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data,omitempty"`
	Error   string          `json:"error,omitempty"`
	Message string          `json:"message,omitempty"`
}

// Session types
type WzapSession struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	APIKey    string    `json:"apiKey,omitempty"`
	JID       string    `json:"jid,omitempty"`
	QRCode    string    `json:"qrCode,omitempty"`
	Connected int       `json:"connected"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type WzapCreateSessionReq struct {
	Name   string `json:"name"`
	APIKey string `json:"apiKey,omitempty"`
}

type WzapCreateSessionResp struct {
	WzapSession
}

// Message types
type WzapSendTextReq struct {
	JID  string `json:"jid"`
	Text string `json:"text"`
}

type WzapSendMediaReq struct {
	JID      string `json:"jid"`
	MimeType string `json:"mimeType,omitempty"`
	Caption  string `json:"caption,omitempty"`
	Filename string `json:"filename,omitempty"`
	Base64   string `json:"base64,omitempty"`
}

type WzapSendContactReq struct {
	JID   string `json:"jid"`
	Name  string `json:"name"`
	Vcard string `json:"vcard"`
}

type WzapSendLocationReq struct {
	JID     string  `json:"jid"`
	Lat     float64 `json:"lat"`
	Lng     float64 `json:"lng"`
	Name    string  `json:"name,omitempty"`
	Address string  `json:"address,omitempty"`
}

type WzapSendPollReq struct {
	JID             string   `json:"jid"`
	Name            string   `json:"name"`
	Options         []string `json:"options"`
	SelectableCount int      `json:"selectableCount"`
}

type WzapSendStickerReq struct {
	JID      string `json:"jid"`
	MimeType string `json:"mimeType"`
	Base64   string `json:"base64"`
}

type WzapSendLinkReq struct {
	JID         string `json:"jid"`
	URL         string `json:"url"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

type WzapEditMessageReq struct {
	JID       string `json:"jid"`
	MessageID string `json:"messageId"`
	Text      string `json:"text"`
}

type WzapDeleteMessageReq struct {
	JID       string `json:"jid"`
	MessageID string `json:"messageId"`
}

type WzapReactMessageReq struct {
	JID       string `json:"jid"`
	MessageID string `json:"messageId"`
	Reaction  string `json:"reaction"`
}

type WzapMarkReadReq struct {
	JID       string `json:"jid"`
	MessageID string `json:"messageId"`
}

type WzapSetPresenceReq struct {
	JID      string `json:"jid"`
	Presence string `json:"presence"`
}

// Group types
type WzapGroupInfo struct {
	JID          string `json:"jid"`
	Name         string `json:"name"`
	Participants int    `json:"participants"`
	IsAdmin      bool   `json:"isAdmin"`
}

type WzapCreateGroupReq struct {
	Name         string   `json:"name"`
	Participants []string `json:"participants"`
}

type WzapGroupJIDReq struct {
	GroupJID string `json:"groupJid"`
}

type WzapGroupJoinReq struct {
	InviteCode string `json:"inviteCode"`
}

type WzapGroupParticipantReq struct {
	GroupJID     string   `json:"groupJid"`
	Participants []string `json:"participants"`
	Action       string   `json:"action"`
}

type WzapGroupRequestActionReq struct {
	GroupJID     string   `json:"groupJid"`
	Participants []string `json:"participants"`
	Action       string   `json:"action"`
}

type WzapGroupTextReq struct {
	GroupJID string `json:"groupJid"`
	Text     string `json:"text"`
}

type WzapGroupPhotoReq struct {
	GroupJID    string `json:"groupJid"`
	PhotoBase64 string `json:"photoBase64"`
}

type WzapGroupSettingReq struct {
	GroupJID string `json:"groupJid"`
	Enabled  bool   `json:"enabled"`
}

// Contact types
type WzapCheckContactReq struct {
	Phones []string `json:"phones"`
}

type WzapCheckContactResp struct {
	Exists      bool   `json:"exists"`
	JID         string `json:"jid,omitempty"`
	PhoneNumber string `json:"phoneNumber"`
}

type WzapContact struct {
	JID      string `json:"jid"`
	Name     string `json:"name,omitempty"`
	PushName string `json:"pushName,omitempty"`
}

type WzapGetAvatarReq struct {
	JID string `json:"jid"`
}

type WzapGetAvatarResp struct {
	URL string `json:"url"`
	ID  string `json:"id"`
}

type WzapBlockContactReq struct {
	JID string `json:"jid"`
}

type WzapGetUserInfoReq struct {
	JIDs []string `json:"jids"`
}

type WzapUserInfoResp struct {
	JID     string   `json:"jid"`
	Status  string   `json:"status"`
	Picture string   `json:"picture"`
	Devices []string `json:"devices"`
}

type WzapSetProfilePictureReq struct {
	Base64 string `json:"base64"`
}

// Label types
type WzapLabelChatReq struct {
	JID     string `json:"jid"`
	LabelID string `json:"labelId"`
}

type WzapLabelMessageReq struct {
	JID       string `json:"jid"`
	LabelID   string `json:"labelId"`
	MessageID string `json:"messageId"`
}

type WzapEditLabelReq struct {
	Color   int32  `json:"color"`
	Deleted bool   `json:"deleted"`
	LabelID string `json:"labelId"`
	Name    string `json:"name"`
}

// Chat types
type WzapChatActionReq struct {
	JID string `json:"jid"`
}

type WzapChatMuteReq struct {
	JID      string `json:"jid"`
	Duration int64  `json:"duration"`
}

// Newsletter types
type WzapCreateNewsletterReq struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Picture     string `json:"picture,omitempty"`
}

type WzapNewsletterMessageReq struct {
	NewsletterJID string `json:"newsletterJid"`
	Count         int    `json:"count"`
	BeforeID      int    `json:"beforeId"`
}

type WzapNewsletterSubscribeReq struct {
	NewsletterJID string `json:"newsletterJid"`
}

// Community types
type WzapCreateCommunityReq struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type WzapCommunityParticipantReq struct {
	CommunityJID string   `json:"communityJid"`
	Participants []string `json:"participants"`
}

// Webhook event from wzap service
type WzapWebhookEvent struct {
	EventID   string `json:"eventId"`
	SessionID string `json:"sessionId"`
	Event     string `json:"event"`
	Timestamp string `json:"timestamp"`

	// Message event fields
	ID          string `json:"id,omitempty"`
	From        string `json:"from,omitempty"`
	PushName    string `json:"pushName,omitempty"`
	Content     string `json:"content,omitempty"`
	ContentType string `json:"contentType,omitempty"`
	FromMe      bool   `json:"fromMe,omitempty"`

	// Receipt event fields
	MessageID string `json:"messageId,omitempty"`
	Status    string `json:"status,omitempty"`

	// Connection event fields
	JID    string `json:"jid,omitempty"`
	Reason int    `json:"reason,omitempty"`

	// QR event fields
	QR string `json:"qr,omitempty"`
}

// Webhook creation for wzap service
type WzapCreateWebhookReq struct {
	URL    string   `json:"url"`
	Events []string `json:"events"`
}
