package domain

// IncomingMessage represents a received WhatsApp message from wzap via webhook.
type IncomingMessage struct {
	ExternalID  string
	From        string
	PushName    string
	Content     string
	MediaURL    string
	MediaType   string
	ContentType ContentType
	IsFromMe    bool
	Timestamp   int64
}

// StatusUpdate represents a message delivery/read receipt from wzap via webhook.
type StatusUpdate struct {
	ExternalID string
	Status     MessageStatus
}
