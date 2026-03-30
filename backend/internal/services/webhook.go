package services

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/felipyfgs/zenwoot/backend/internal/logger"
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

func (s *WebhookService) GetByID(ctx context.Context, accountID, id int64) (*models.Webhook, error) {
	return s.webhookRepo.GetByID(ctx, accountID, id)
}

func (s *WebhookService) Update(ctx context.Context, m *models.Webhook) (*models.Webhook, error) {
	if err := s.webhookRepo.Update(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *WebhookService) Dispatch(ctx context.Context, accountID int64, event string, payload any) {
	hooks, err := s.webhookRepo.ListActiveByEvent(ctx, accountID, event)
	if err != nil {
		logger.Error().Err(err).Msg("failed to list webhooks for dispatch")
		return
	}
	if len(hooks) == 0 {
		return
	}

	body, err := json.Marshal(payload)
	if err != nil {
		logger.Error().Err(err).Msg("failed to marshal webhook payload")
		return
	}

	for _, wh := range hooks {
		go s.deliver(wh, body)
	}
}

func (s *WebhookService) deliver(wh *models.Webhook, body []byte) {
	req, err := http.NewRequest(http.MethodPost, wh.URL, bytes.NewReader(body))
	if err != nil {
		logger.Error().Err(err).Str("url", wh.URL).Msg("failed to create webhook request")
		return
	}
	req.Header.Set("Content-Type", "application/json")
	if wh.HmacToken != nil {
		req.Header.Set("X-Zenwoot-Signature", s.sign(*wh.HmacToken, body))
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		logger.Error().Err(err).Str("url", wh.URL).Msg("failed to deliver webhook")
		return
	}
	defer resp.Body.Close()

	if _, err := io.Copy(io.Discard, resp.Body); err != nil {
		logger.Error().Err(err).Str("url", wh.URL).Msg("failed to read webhook response body")
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		logger.Warn().Int("status", resp.StatusCode).Str("url", wh.URL).Msg("webhook returned non-success status")
	}
}

func (s *WebhookService) sign(secret string, body []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}
