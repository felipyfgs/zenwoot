package workers

import (
	"context"
	"encoding/json"

	"github.com/nats-io/nats.go"

	"github.com/felipyfgs/zenwoot/backend/internal/logger"
	"github.com/felipyfgs/zenwoot/backend/internal/repo"
)

type AutomationWorker struct {
	nc       *nats.Conn
	autoRepo *repo.AutomationRuleRepo
	subs     []*nats.Subscription
}

func NewAutomationWorker(nc *nats.Conn, autoRepo *repo.AutomationRuleRepo) *AutomationWorker {
	return &AutomationWorker{nc: nc, autoRepo: autoRepo}
}

func (w *AutomationWorker) Start() error {
	subjects := []string{
		"zenwoot.conversation.created",
		"zenwoot.message.created",
		"zenwoot.conversation.resolved",
	}
	for _, subj := range subjects {
		eventName := subjectToEventName(subj)
		sub, err := w.nc.Subscribe(subj, w.handle(eventName))
		if err != nil {
			return err
		}
		w.subs = append(w.subs, sub)
	}
	logger.Info().Msg("AutomationWorker started")
	return nil
}

func (w *AutomationWorker) Stop() {
	for _, sub := range w.subs {
		_ = sub.Unsubscribe()
	}
}

func (w *AutomationWorker) handle(eventName string) nats.MsgHandler {
	return func(msg *nats.Msg) {
		if w.autoRepo == nil {
			logger.Error().Msg("autoRepo is nil in AutomationWorker")
			return
		}

		var payload map[string]any
		if err := json.Unmarshal(msg.Data, &payload); err != nil {
			logger.Warn().Err(err).Msg("failed to unmarshal automation event")
			return
		}

		accountIDf, ok := payload["account_id"].(float64)
		if !ok {
			logger.Warn().Msg("account_id not found in payload or invalid type")
			return
		}
		accountID := int64(accountIDf)
		if accountID == 0 {
			return
		}

		rules, err := w.autoRepo.ListActive(context.Background(), accountID, eventName)
		if err != nil {
			logger.Error().Err(err).Int64("account_id", accountID).Str("event", eventName).Msg("failed to fetch automation rules")
			return
		}
		if len(rules) == 0 {
			return
		}

		for _, rule := range rules {
			w.executeRule(rule)
		}
	}
}

func (w *AutomationWorker) executeRule(rule any) {
	logger.Debug().Interface("rule", rule).Msg("executing automation rule")
}

func subjectToEventName(subject string) string {
	mapping := map[string]string{
		"zenwoot.conversation.created":  "conversation_created",
		"zenwoot.conversation.updated":  "conversation_updated",
		"zenwoot.conversation.resolved": "conversation_resolved",
		"zenwoot.message.created":       "message_created",
	}
	if v, ok := mapping[subject]; ok {
		return v
	}
	return subject
}
