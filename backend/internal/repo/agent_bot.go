package repo

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"

	"github.com/felipyfgs/zenwoot/backend/internal/models"
)

type AgentBotRepo struct {
	BaseRepo[models.AgentBot]
}

func NewAgentBotRepo(db *bun.DB) *AgentBotRepo {
	return &AgentBotRepo{BaseRepo: *NewBaseRepo[models.AgentBot](db)}
}

func (r *AgentBotRepo) ListByAccount(ctx context.Context, accountID int64) ([]*models.AgentBot, error) {
	var items []*models.AgentBot
	err := r.WithTenant(ctx, accountID).OrderExpr(`"name" ASC`).Scan(ctx, &items)
	if err != nil {
		return nil, fmt.Errorf("agentBotRepo.ListByAccount: %w", err)
	}
	return items, nil
}

func (r *AgentBotRepo) Create(ctx context.Context, m *models.AgentBot) error {
	_, err := r.DB().NewInsert().Model(m).Exec(ctx)
	if err != nil {
		return fmt.Errorf("agentBotRepo.Create: %w", err)
	}
	return nil
}

func (r *AgentBotRepo) Update(ctx context.Context, m *models.AgentBot) error {
	_, err := r.DB().NewUpdate().Model(m).
		Where(`"id" = ? AND "account_id" = ?`, m.ID, m.AccountID).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("agentBotRepo.Update: %w", err)
	}
	return nil
}

func (r *AgentBotRepo) Delete(ctx context.Context, accountID, id int64) error {
	_, err := r.DB().NewDelete().TableExpr(`"agent_bots"`).
		Where(`"id" = ? AND "account_id" = ?`, id, accountID).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("agentBotRepo.Delete: %w", err)
	}
	return nil
}

func (r *AgentBotRepo) GetByID(ctx context.Context, accountID, id int64) (*models.AgentBot, error) {
	var m models.AgentBot
	err := r.WithTenant(ctx, accountID).Where(`"id" = ?`, id).Scan(ctx, &m)
	if err != nil {
		return nil, fmt.Errorf("agentBotRepo.GetByID: %w", err)
	}
	return &m, nil
}
