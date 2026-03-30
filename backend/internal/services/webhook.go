package services

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"github.com/felipyfgs/zenwoot/backend/internal/models"
	"github.com/felipyfgs/zenwoot/backend/internal/repo"
)

type WebhookService struct {
	webhookRepo *repo.WebhookRepo
	httpClient  *http.Client
}

func NewWebhookService(webhookRepo *repo.WebhookRepo) *WebhookService {
	return &WebhookService{
		webhookRepo: webhookRepo,
		httpClient:  &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *WebhookService) List(ctx context.Context, accountID int64) ([]*models.Webhook, error) {
	return s.webhookRepo.ListByAccount(ctx, accountID)
}

func (s *WebhookService) Create(ctx context.Context, m *models.Webhook) (*models.Webhook, error) {
	if err := s.webhookRepo.Create(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *WebhookService) Delete(ctx context.Context, accountID, id int64) error {
	return s.webhookRepo.Delete(ctx, accountID, id)
}

func (s *WebhookService) Dispatch(ctx context.Context, accountID int64, event string, payload any) {
	hooks, err := s.webhookRepo.ListActiveByEvent(ctx, accountID, event)
	if err != nil {
		return
	}
	body, _ := json.Marshal(payload)
	for _, wh := range hooks {
		go s.deliver(wh, body)
	}
}

func (s *WebhookService) deliver(wh *models.Webhook, body []byte) {
	req, err := http.NewRequest(http.MethodPost, wh.URL, bytes.NewReader(body))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	if wh.HmacToken != nil {
		req.Header.Set("X-Zenwoot-Signature", s.sign(*wh.HmacToken, body))
	}
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
}

func (s *WebhookService) sign(secret string, body []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}
