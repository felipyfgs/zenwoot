package ws

import (
	"encoding/json"
	"sync"

	"github.com/rs/zerolog/log"
)

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type Client struct {
	accountID string
	inboxIDs  map[string]bool
	userID    string      // Optional: for typing indicators
	Send      chan []byte // Exported for WebSocket handler
	hub       *Hub
}

type Hub struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan *broadcastMsg
	mu         sync.RWMutex
}

type broadcastMsg struct {
	accountID string
	inboxID   string
	data      []byte
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client, 16),
		unregister: make(chan *Client, 16),
		broadcast:  make(chan *broadcastMsg, 256),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.register:
			h.mu.Lock()
			h.clients[c] = true
			h.mu.Unlock()

		case c := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				close(c.Send)
			}
			h.mu.Unlock()

		case msg := <-h.broadcast:
			h.mu.RLock()
			for c := range h.clients {
				if c.accountID != msg.accountID {
					continue
				}
				if msg.inboxID != "" && !c.inboxIDs[msg.inboxID] {
					continue
				}
				select {
				case c.Send <- msg.data:
				default:
					log.Warn().Str("accountId", c.accountID).Msg("WS client send buffer full, dropping message")
				}
			}
			h.mu.RUnlock()
		}
	}
}

// Publish sends an event to all connected WS clients subscribed to an account/inbox.
func (h *Hub) Publish(accountID, inboxID string, eventType string, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal WS payload")
		return
	}

	evt := Event{Type: eventType, Payload: json.RawMessage(data)}
	raw, _ := json.Marshal(evt)

	h.broadcast <- &broadcastMsg{
		accountID: accountID,
		inboxID:   inboxID,
		data:      raw,
	}
}

func (h *Hub) NewClient(accountID string, inboxIDs []string) *Client {
	ids := make(map[string]bool, len(inboxIDs))
	for _, id := range inboxIDs {
		ids[id] = true
	}
	c := &Client{
		accountID: accountID,
		inboxIDs:  ids,
		Send:      make(chan []byte, 64),
		hub:       h,
	}
	h.register <- c
	return c
}

func (h *Hub) NewClientWithUser(accountID string, inboxIDs []string, userID string) *Client {
	ids := make(map[string]bool, len(inboxIDs))
	for _, id := range inboxIDs {
		ids[id] = true
	}
	c := &Client{
		accountID: accountID,
		inboxIDs:  ids,
		userID:    userID,
		Send:      make(chan []byte, 64),
		hub:       h,
	}
	h.register <- c
	return c
}

func (h *Hub) RemoveClient(c *Client) {
	h.unregister <- c
}

// TypingEvent represents a typing indicator event
type TypingEvent struct {
	ConversationID string `json:"conversationId"`
	UserID         string `json:"userId"`
	UserName       string `json:"userName,omitempty"`
	IsTyping       bool   `json:"isTyping"`
}

// PublishTyping sends a typing indicator event to all connected WS clients
func (h *Hub) PublishTyping(accountID, inboxID, conversationID, userID, userName string, isTyping bool) {
	evt := TypingEvent{
		ConversationID: conversationID,
		UserID:         userID,
		UserName:       userName,
		IsTyping:       isTyping,
	}
	h.Publish(accountID, inboxID, "typing", evt)
}
