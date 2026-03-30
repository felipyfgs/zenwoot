package whatsapp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"wzap/internal/domain"
)

// WzapClient is an HTTP client that calls the wzap API as an external WhatsApp provider.
// It replaces the older embedded approach.
type WzapClient struct {
	baseURL    string
	adminKey   string
	httpClient *http.Client

	// Callbacks registered by the engine
	msgHandler          func(channelID string, msg *domain.IncomingMessage)
	statusHandler       func(channelID string, status *domain.StatusUpdate)
	connectedHandler    func(channelID string, jid string)
	disconnectedHandler func(channelID string)
	loggedOutHandler    func(channelID string)
	qrHandler           func(channelID string, qr string)
}

func NewWzapClient(baseURL, adminKey string) *WzapClient {
	return &WzapClient{
		baseURL:  strings.TrimRight(baseURL, "/"),
		adminKey: adminKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ─── Session Management ────────────────────────────────────────────────────────

// CreateSession creates a new session in wzap for the given channelID.
// Returns the wzap session ID (which equals the channelID we pass as Name).
func (c *WzapClient) CreateSession(ctx context.Context, channelID string) (*WzapSession, error) {
	req := WzapCreateSessionReq{Name: channelID}
	var resp WzapCreateSessionResp
	if err := c.post(ctx, "/sessions", req, &resp, c.adminKey); err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}
	return &resp.WzapSession, nil
}

// GetSession reads a session from wzap by channelID (used as the session name/ID).
func (c *WzapClient) GetSession(ctx context.Context, channelID string) (*WzapSession, error) {
	var resp WzapSession
	if err := c.get(ctx, "/sessions/"+channelID, &resp, c.adminKey); err != nil {
		return nil, fmt.Errorf("get session: %w", err)
	}
	return &resp, nil
}

// Connect starts connection / QR pairing for a session.
func (c *WzapClient) Connect(ctx context.Context, channelID string, _ string) (interface{}, error) {
	if err := c.postEmpty(ctx, "/sessions/"+channelID+"/connect", c.adminKey); err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}
	return nil, nil
}

// Disconnect disconnects a session in wzap.
func (c *WzapClient) Disconnect(channelID string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = c.postEmpty(ctx, "/sessions/"+channelID+"/disconnect", c.adminKey)
}

// IsConnected checks if a session is connected by querying wzap.
func (c *WzapClient) IsConnected(channelID string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	sess, err := c.GetSession(ctx, channelID)
	if err != nil {
		return false
	}
	return sess.Connected == 1
}

// ReconnectAll reconnects all sessions that have a stored JID in wzap.
// The resolver maps JID -> channelID (not needed for wzap, sessions are self-aware).
func (c *WzapClient) ReconnectAll(_ context.Context, _ func(jid string) (string, error)) error {
	// wzap manages its own sessions; on startup it auto-reconnects.
	return nil
}

// ─── Messaging ────────────────────────────────────────────────────────────────

func (c *WzapClient) SendText(ctx context.Context, channelID, apiKey, to, text string) (string, error) {
	req := WzapSendTextReq{JID: toJID(to), Text: text}
	var resp struct {
		MessageID string `json:"messageId"`
	}
	if err := c.post(ctx, "/sessions/"+channelID+"/messages/text", req, &resp, apiKey); err != nil {
		return "", fmt.Errorf("send text: %w", err)
	}
	return resp.MessageID, nil
}

func (c *WzapClient) SendMedia(ctx context.Context, channelID, apiKey string, req WzapSendMediaReq) (string, error) {
	req.JID = toJID(req.JID)
	var resp struct {
		MessageID string `json:"messageId"`
	}
	if err := c.post(ctx, "/sessions/"+channelID+"/messages/image", req, &resp, apiKey); err != nil {
		return "", fmt.Errorf("send media: %w", err)
	}
	return resp.MessageID, nil
}

func (c *WzapClient) SendContact(ctx context.Context, channelID, apiKey string, req WzapSendContactReq) (string, error) {
	req.JID = toJID(req.JID)
	var resp struct{ MessageID string `json:"messageId"` }
	if err := c.post(ctx, "/sessions/"+channelID+"/messages/contact", req, &resp, apiKey); err != nil {
		return "", fmt.Errorf("send contact: %w", err)
	}
	return resp.MessageID, nil
}

func (c *WzapClient) SendLocation(ctx context.Context, channelID, apiKey string, req WzapSendLocationReq) (string, error) {
	req.JID = toJID(req.JID)
	var resp struct{ MessageID string `json:"messageId"` }
	if err := c.post(ctx, "/sessions/"+channelID+"/messages/location", req, &resp, apiKey); err != nil {
		return "", fmt.Errorf("send location: %w", err)
	}
	return resp.MessageID, nil
}

func (c *WzapClient) SendPoll(ctx context.Context, channelID, apiKey string, req WzapSendPollReq) (string, error) {
	req.JID = toJID(req.JID)
	var resp struct{ MessageID string `json:"messageId"` }
	if err := c.post(ctx, "/sessions/"+channelID+"/messages/poll", req, &resp, apiKey); err != nil {
		return "", fmt.Errorf("send poll: %w", err)
	}
	return resp.MessageID, nil
}

func (c *WzapClient) SendSticker(ctx context.Context, channelID, apiKey string, req WzapSendStickerReq) (string, error) {
	req.JID = toJID(req.JID)
	var resp struct{ MessageID string `json:"messageId"` }
	if err := c.post(ctx, "/sessions/"+channelID+"/messages/sticker", req, &resp, apiKey); err != nil {
		return "", fmt.Errorf("send sticker: %w", err)
	}
	return resp.MessageID, nil
}

func (c *WzapClient) SendLink(ctx context.Context, channelID, apiKey string, req WzapSendLinkReq) (string, error) {
	req.JID = toJID(req.JID)
	var resp struct{ MessageID string `json:"messageId"` }
	if err := c.post(ctx, "/sessions/"+channelID+"/messages/link", req, &resp, apiKey); err != nil {
		return "", fmt.Errorf("send link: %w", err)
	}
	return resp.MessageID, nil
}

func (c *WzapClient) EditMessage(ctx context.Context, channelID, apiKey string, req WzapEditMessageReq) (string, error) {
	req.JID = toJID(req.JID)
	var resp struct{ MessageID string `json:"messageId"` }
	if err := c.post(ctx, "/sessions/"+channelID+"/messages/edit", req, &resp, apiKey); err != nil {
		return "", fmt.Errorf("edit message: %w", err)
	}
	return resp.MessageID, nil
}

func (c *WzapClient) DeleteMessage(ctx context.Context, channelID, apiKey string, req WzapDeleteMessageReq) (string, error) {
	req.JID = toJID(req.JID)
	var resp struct{ MessageID string `json:"messageId"` }
	if err := c.post(ctx, "/sessions/"+channelID+"/messages/delete", req, &resp, apiKey); err != nil {
		return "", fmt.Errorf("delete message: %w", err)
	}
	return resp.MessageID, nil
}

func (c *WzapClient) ReactMessage(ctx context.Context, channelID, apiKey string, req WzapReactMessageReq) (string, error) {
	req.JID = toJID(req.JID)
	var resp struct{ MessageID string `json:"messageId"` }
	if err := c.post(ctx, "/sessions/"+channelID+"/messages/reaction", req, &resp, apiKey); err != nil {
		return "", fmt.Errorf("react message: %w", err)
	}
	return resp.MessageID, nil
}

func (c *WzapClient) MarkRead(ctx context.Context, channelID, apiKey string, req WzapMarkReadReq) error {
	req.JID = toJID(req.JID)
	return c.post(ctx, "/sessions/"+channelID+"/messages/read", req, nil, apiKey)
}

func (c *WzapClient) SetPresence(ctx context.Context, channelID, apiKey string, req WzapSetPresenceReq) error {
	req.JID = toJID(req.JID)
	return c.post(ctx, "/sessions/"+channelID+"/messages/presence", req, nil, apiKey)
}

// ─── Contacts ─────────────────────────────────────────────────────────────────

func (c *WzapClient) CheckContacts(ctx context.Context, channelID, apiKey string, phones []string) ([]WzapCheckContactResp, error) {
	req := WzapCheckContactReq{Phones: phones}
	var resp []WzapCheckContactResp
	if err := c.post(ctx, "/sessions/"+channelID+"/contacts/check", req, &resp, apiKey); err != nil {
		return nil, fmt.Errorf("check contacts: %w", err)
	}
	return resp, nil
}

func (c *WzapClient) ListContacts(ctx context.Context, channelID, apiKey string) ([]WzapContact, error) {
	var resp []WzapContact
	if err := c.get(ctx, "/sessions/"+channelID+"/contacts", &resp, apiKey); err != nil {
		return nil, fmt.Errorf("list contacts: %w", err)
	}
	return resp, nil
}

func (c *WzapClient) GetAvatar(ctx context.Context, channelID, apiKey string, req WzapGetAvatarReq) (*WzapGetAvatarResp, error) {
	var resp WzapGetAvatarResp
	if err := c.post(ctx, "/sessions/"+channelID+"/contacts/avatar", req, &resp, apiKey); err != nil {
		return nil, fmt.Errorf("get avatar: %w", err)
	}
	return &resp, nil
}

func (c *WzapClient) BlockContact(ctx context.Context, channelID, apiKey string, req WzapBlockContactReq) error {
	return c.post(ctx, "/sessions/"+channelID+"/contacts/block", req, nil, apiKey)
}

func (c *WzapClient) UnblockContact(ctx context.Context, channelID, apiKey string, req WzapBlockContactReq) error {
	return c.post(ctx, "/sessions/"+channelID+"/contacts/unblock", req, nil, apiKey)
}

func (c *WzapClient) GetBlocklist(ctx context.Context, channelID, apiKey string) (interface{}, error) {
	var resp interface{}
	if err := c.get(ctx, "/sessions/"+channelID+"/contacts/blocklist", &resp, apiKey); err != nil {
		return nil, fmt.Errorf("get blocklist: %w", err)
	}
	return resp, nil
}

func (c *WzapClient) GetUserInfo(ctx context.Context, channelID, apiKey string, req WzapGetUserInfoReq) ([]WzapUserInfoResp, error) {
	var resp []WzapUserInfoResp
	if err := c.post(ctx, "/sessions/"+channelID+"/contacts/info", req, &resp, apiKey); err != nil {
		return nil, fmt.Errorf("get user info: %w", err)
	}
	return resp, nil
}

func (c *WzapClient) GetPrivacySettings(ctx context.Context, channelID, apiKey string) (interface{}, error) {
	var resp interface{}
	if err := c.get(ctx, "/sessions/"+channelID+"/contacts/privacy", &resp, apiKey); err != nil {
		return nil, fmt.Errorf("get privacy settings: %w", err)
	}
	return resp, nil
}

func (c *WzapClient) SetProfilePicture(ctx context.Context, channelID, apiKey string, req WzapSetProfilePictureReq) (string, error) {
	var resp struct{ PictureID string `json:"pictureId"` }
	if err := c.post(ctx, "/sessions/"+channelID+"/contacts/profile-picture", req, &resp, apiKey); err != nil {
		return "", fmt.Errorf("set profile picture: %w", err)
	}
	return resp.PictureID, nil
}

// ─── Groups ───────────────────────────────────────────────────────────────────

func (c *WzapClient) ListGroups(ctx context.Context, channelID, apiKey string) ([]WzapGroupInfo, error) {
	var resp []WzapGroupInfo
	if err := c.get(ctx, "/sessions/"+channelID+"/groups", &resp, apiKey); err != nil {
		return nil, fmt.Errorf("list groups: %w", err)
	}
	return resp, nil
}

func (c *WzapClient) CreateGroup(ctx context.Context, channelID, apiKey string, req WzapCreateGroupReq) (*WzapGroupInfo, error) {
	var resp WzapGroupInfo
	if err := c.post(ctx, "/sessions/"+channelID+"/groups/create", req, &resp, apiKey); err != nil {
		return nil, fmt.Errorf("create group: %w", err)
	}
	return &resp, nil
}

func (c *WzapClient) GetGroupInfo(ctx context.Context, channelID, apiKey, groupJID string) (*WzapGroupInfo, error) {
	var resp WzapGroupInfo
	if err := c.post(ctx, "/sessions/"+channelID+"/groups/info", WzapGroupJIDReq{GroupJID: groupJID}, &resp, apiKey); err != nil {
		return nil, fmt.Errorf("get group info: %w", err)
	}
	return &resp, nil
}

func (c *WzapClient) GetGroupInfoFromLink(ctx context.Context, channelID, apiKey, inviteCode string) (*WzapGroupInfo, error) {
	var resp WzapGroupInfo
	if err := c.post(ctx, "/sessions/"+channelID+"/groups/invite-info", WzapGroupJoinReq{InviteCode: inviteCode}, &resp, apiKey); err != nil {
		return nil, fmt.Errorf("get group info from link: %w", err)
	}
	return &resp, nil
}

func (c *WzapClient) JoinGroupWithLink(ctx context.Context, channelID, apiKey, inviteCode string) (string, error) {
	var resp struct{ JID string `json:"jid"` }
	if err := c.post(ctx, "/sessions/"+channelID+"/groups/join", WzapGroupJoinReq{InviteCode: inviteCode}, &resp, apiKey); err != nil {
		return "", fmt.Errorf("join group: %w", err)
	}
	return resp.JID, nil
}

func (c *WzapClient) GetGroupInviteLink(ctx context.Context, channelID, apiKey, groupJID string) (string, error) {
	var resp struct{ InviteLink string `json:"inviteLink"` }
	if err := c.post(ctx, "/sessions/"+channelID+"/groups/invite-link", WzapGroupJIDReq{GroupJID: groupJID}, &resp, apiKey); err != nil {
		return "", fmt.Errorf("get invite link: %w", err)
	}
	return resp.InviteLink, nil
}

func (c *WzapClient) LeaveGroup(ctx context.Context, channelID, apiKey, groupJID string) error {
	return c.post(ctx, "/sessions/"+channelID+"/groups/leave", WzapGroupJIDReq{GroupJID: groupJID}, nil, apiKey)
}

func (c *WzapClient) UpdateGroupParticipants(ctx context.Context, channelID, apiKey string, req WzapGroupParticipantReq) (interface{}, error) {
	var resp interface{}
	if err := c.post(ctx, "/sessions/"+channelID+"/groups/participants", req, &resp, apiKey); err != nil {
		return nil, fmt.Errorf("update participants: %w", err)
	}
	return resp, nil
}

func (c *WzapClient) GetGroupRequests(ctx context.Context, channelID, apiKey, groupJID string) (interface{}, error) {
	var resp interface{}
	if err := c.post(ctx, "/sessions/"+channelID+"/groups/requests", WzapGroupJIDReq{GroupJID: groupJID}, &resp, apiKey); err != nil {
		return nil, fmt.Errorf("get group requests: %w", err)
	}
	return resp, nil
}

func (c *WzapClient) UpdateGroupRequests(ctx context.Context, channelID, apiKey string, req WzapGroupRequestActionReq) (interface{}, error) {
	var resp interface{}
	if err := c.post(ctx, "/sessions/"+channelID+"/groups/requests/action", req, &resp, apiKey); err != nil {
		return nil, fmt.Errorf("update group requests: %w", err)
	}
	return resp, nil
}

func (c *WzapClient) UpdateGroupName(ctx context.Context, channelID, apiKey string, req WzapGroupTextReq) error {
	return c.post(ctx, "/sessions/"+channelID+"/groups/name", req, nil, apiKey)
}

func (c *WzapClient) UpdateGroupDescription(ctx context.Context, channelID, apiKey string, req WzapGroupTextReq) error {
	return c.post(ctx, "/sessions/"+channelID+"/groups/description", req, nil, apiKey)
}

func (c *WzapClient) UpdateGroupPhoto(ctx context.Context, channelID, apiKey string, req WzapGroupPhotoReq) (string, error) {
	var resp struct{ PictureID string `json:"pictureId"` }
	if err := c.post(ctx, "/sessions/"+channelID+"/groups/photo", req, &resp, apiKey); err != nil {
		return "", fmt.Errorf("update group photo: %w", err)
	}
	return resp.PictureID, nil
}

func (c *WzapClient) SetGroupAnnounce(ctx context.Context, channelID, apiKey string, req WzapGroupSettingReq) error {
	return c.post(ctx, "/sessions/"+channelID+"/groups/announce", req, nil, apiKey)
}

func (c *WzapClient) SetGroupLocked(ctx context.Context, channelID, apiKey string, req WzapGroupSettingReq) error {
	return c.post(ctx, "/sessions/"+channelID+"/groups/locked", req, nil, apiKey)
}

func (c *WzapClient) SetGroupJoinApproval(ctx context.Context, channelID, apiKey string, req WzapGroupSettingReq) error {
	return c.post(ctx, "/sessions/"+channelID+"/groups/join-approval", req, nil, apiKey)
}

// ─── Chat ─────────────────────────────────────────────────────────────────────

func (c *WzapClient) ChatAction(ctx context.Context, channelID, apiKey, action string, req interface{}) error {
	return c.post(ctx, "/sessions/"+channelID+"/chat/"+action, req, nil, apiKey)
}

// ─── Labels ───────────────────────────────────────────────────────────────────

func (c *WzapClient) LabelChat(ctx context.Context, channelID, apiKey string, req WzapLabelChatReq, add bool) error {
	path := "/sessions/" + channelID + "/label/chat"
	if !add {
		path = "/sessions/" + channelID + "/unlabel/chat"
	}
	return c.post(ctx, path, req, nil, apiKey)
}

func (c *WzapClient) LabelMessage(ctx context.Context, channelID, apiKey string, req WzapLabelMessageReq, add bool) error {
	path := "/sessions/" + channelID + "/label/message"
	if !add {
		path = "/sessions/" + channelID + "/unlabel/message"
	}
	return c.post(ctx, path, req, nil, apiKey)
}

func (c *WzapClient) EditLabel(ctx context.Context, channelID, apiKey string, req WzapEditLabelReq) error {
	return c.post(ctx, "/sessions/"+channelID+"/label/edit", req, nil, apiKey)
}

// ─── Newsletter ────────────────────────────────────────────────────────────────

func (c *WzapClient) NewsletterAction(ctx context.Context, channelID, apiKey, path string, req interface{}, out interface{}) error {
	return c.post(ctx, "/sessions/"+channelID+"/newsletter/"+path, req, out, apiKey)
}

func (c *WzapClient) ListNewsletters(ctx context.Context, channelID, apiKey string) (interface{}, error) {
	var resp interface{}
	if err := c.get(ctx, "/sessions/"+channelID+"/newsletter/list", &resp, apiKey); err != nil {
		return nil, err
	}
	return resp, nil
}

// ─── Community ────────────────────────────────────────────────────────────────

func (c *WzapClient) CommunityAction(ctx context.Context, channelID, apiKey, path string, req interface{}) error {
	return c.post(ctx, "/sessions/"+channelID+"/community/"+path, req, nil, apiKey)
}

// ─── Webhooks ─────────────────────────────────────────────────────────────────

// RegisterWebhook registers a webhook in wzap for the given session.
func (c *WzapClient) RegisterWebhook(ctx context.Context, channelID string, webhookURL string, events []string) error {
	req := WzapCreateWebhookReq{URL: webhookURL, Events: events}
	return c.post(ctx, "/sessions/"+channelID+"/webhooks", req, nil, c.adminKey)
}

// ─── Callback registration ────────────────────────────────────────────────────
// These allow the engine to register handlers for events received via webhook.

func (c *WzapClient) OnMessage(h func(channelID string, msg *domain.IncomingMessage))         { c.msgHandler = h }
func (c *WzapClient) OnStatus(h func(channelID string, s *domain.StatusUpdate))               { c.statusHandler = h }
func (c *WzapClient) OnConnected(h func(channelID string, jid string))                 { c.connectedHandler = h }
func (c *WzapClient) OnDisconnected(h func(channelID string))                          { c.disconnectedHandler = h }
func (c *WzapClient) OnLoggedOut(h func(channelID string))                             { c.loggedOutHandler = h }
func (c *WzapClient) OnQRCode(h func(channelID string, qr string))                     { c.qrHandler = h }

// DispatchWebhookEvent dispatches an incoming event from wzap webhook to the correct handler.
func (c *WzapClient) DispatchWebhookEvent(channelID string, evt *WzapWebhookEvent) {
	switch evt.Event {
	case "message":
		if c.msgHandler != nil {
			contentType := domain.ContentText
			if evt.ContentType != "" {
				contentType = domain.ContentType(evt.ContentType)
			}
			c.msgHandler(channelID, &domain.IncomingMessage{
				ExternalID:  evt.ID,
				From:        evt.From,
				PushName:    evt.PushName,
				Content:     evt.Content,
				ContentType: contentType,
				IsFromMe:    evt.FromMe,
			})
		}
	case "message.ack", "receipt":
		if c.statusHandler != nil {
			c.statusHandler(channelID, &domain.StatusUpdate{
				ExternalID: evt.MessageID,
				Status:     domain.MessageStatus(evt.Status),
			})
		}
	case "connection.update":
		if evt.JID != "" && c.connectedHandler != nil {
			c.connectedHandler(channelID, evt.JID)
		} else if c.disconnectedHandler != nil {
			c.disconnectedHandler(channelID)
		}
	case "qr":
		if c.qrHandler != nil {
			c.qrHandler(channelID, evt.QR)
		}
	case "logout":
		if c.loggedOutHandler != nil {
			c.loggedOutHandler(channelID)
		}
	}
}

// ─── HTTP helpers ──────────────────────────────────────────────────────────────

func (c *WzapClient) post(ctx context.Context, path string, body, out interface{}, apiKey string) error {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("ApiKey", apiKey)

	return c.do(req, out)
}

func (c *WzapClient) postEmpty(ctx context.Context, path, apiKey string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, nil)
	if err != nil {
		return err
	}
	req.Header.Set("ApiKey", apiKey)
	return c.do(req, nil)
}

func (c *WzapClient) get(ctx context.Context, path string, out interface{}, apiKey string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return err
	}
	req.Header.Set("ApiKey", apiKey)
	return c.do(req, out)
}

func (c *WzapClient) do(req *http.Request, out interface{}) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 400 {
		return fmt.Errorf("wzap API error %d: %s", resp.StatusCode, string(body))
	}

	if out == nil {
		return nil
	}

	// Unwrap the wzapAPIResponse envelope
	var envelope wzapAPIResponse
	if err := json.Unmarshal(body, &envelope); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}
	if !envelope.Success {
		msg := envelope.Error
		if msg == "" {
			msg = envelope.Message
		}
		return fmt.Errorf("wzap error: %s", msg)
	}
	if len(envelope.Data) > 0 {
		return json.Unmarshal(envelope.Data, out)
	}
	return nil
}

// toJID normalises a phone number to a WhatsApp JID string.
func toJID(phone string) string {
	phone = strings.TrimSpace(phone)
	if phone == "" {
		return ""
	}
	if strings.Contains(phone, "@") {
		return phone
	}
	return phone + "@s.whatsapp.net"
}
