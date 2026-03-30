package workers

import (
	"context"
	"time"

	"github.com/uptrace/bun"

	"github.com/felipyfgs/zenwoot/backend/internal/logger"
	"github.com/felipyfgs/zenwoot/backend/internal/models"
)

type AutoResolveWorker struct {
	db     *bun.DB
	ticker *time.Ticker
	done   chan struct{}
}

func NewAutoResolveWorker(db *bun.DB, interval time.Duration) *AutoResolveWorker {
	return &AutoResolveWorker{
		db:     db,
		ticker: time.NewTicker(interval),
		done:   make(chan struct{}),
	}
}

func (w *AutoResolveWorker) Start() {
	go func() {
		for {
			select {
			case <-w.ticker.C:
				w.process()
			case <-w.done:
				return
			}
		}
	}()
	logger.Info().Msg("AutoResolveWorker started")
}

func (w *AutoResolveWorker) Stop() {
	w.ticker.Stop()
	close(w.done)
}

func (w *AutoResolveWorker) process() {
	if w.db == nil {
		logger.Error().Msg("db is nil in AutoResolveWorker")
		return
	}

	ctx := context.Background()
	_, err := w.db.NewUpdate().
		TableExpr(`"conversations"`).
		Set(`"status" = ?, "snoozed_until" = NULL`, models.ConvStatusOpen).
		Where(`"status" = ? AND "snoozed_until" <= ?`, models.ConvStatusSnoozed, time.Now()).
		Exec(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("AutoResolveWorker.process error")
	}
}
