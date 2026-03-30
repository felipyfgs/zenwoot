package dto

import "wzap/internal/model"

type WebhookCreateInline struct {
	URL    string            `json:"url"`
	Events []model.EventType `json:"events,omitempty"`
}

type InboxCreateReq struct {
	Name     string                 `json:"name"`
	Webhook  *WebhookCreateInline   `json:"webhook,omitempty"`
	Settings map[string]interface{} `json:"settings,omitempty"`
}

type InboxCreatedResp struct {
	InboxID string         `json:"inboxId"`
	Name    string         `json:"name"`
	Status  string         `json:"status"`
	Webhook *model.Webhook `json:"webhook,omitempty"`
}
