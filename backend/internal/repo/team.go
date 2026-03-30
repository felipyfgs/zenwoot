package repo

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"

	"github.com/felipyfgs/zenwoot/backend/internal/models"
)

type TeamRepo struct {
	BaseRepo[models.Team]
}

func NewTeamRepo(db *bun.DB) *TeamRepo {
	return &TeamRepo{BaseRepo: *NewBaseRepo[models.Team](db)}
}

func (r *TeamRepo) GetByID(ctx context.Context, accountID, id int64) (*models.Team, error) {
	var m models.Team
	err := r.WithTenant(ctx, accountID).Where(`"id" = ?`, id).Scan(ctx, &m)
	if err != nil {
		return nil, fmt.Errorf("teamRepo.GetByID: %w", err)
	}
	return &m, nil
}

func (r *TeamRepo) ListByAccount(ctx context.Context, accountID int64) ([]*models.Team, error) {
	var items []*models.Team
	err := r.WithTenant(ctx, accountID).OrderExpr(`"name" ASC`).Scan(ctx, &items)
	if err != nil {
		return nil, fmt.Errorf("teamRepo.ListByAccount: %w", err)
	}
	return items, nil
}

func (r *TeamRepo) ListMembers(ctx context.Context, teamID int64) ([]*models.TeamMember, error) {
	var items []*models.TeamMember
	err := r.DB().NewSelect().Model(&items).
		Where(`"team_id" = ?`, teamID).
		Relation("User").
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("teamRepo.ListMembers: %w", err)
	}
	return items, nil
}

func (r *TeamRepo) Create(ctx context.Context, m *models.Team) error {
	_, err := r.DB().NewInsert().Model(m).Exec(ctx)
	if err != nil {
		return fmt.Errorf("teamRepo.Create: %w", err)
	}
	return nil
}

func (r *TeamRepo) Update(ctx context.Context, m *models.Team) error {
	_, err := r.DB().NewUpdate().Model(m).
		Where(`"id" = ? AND "account_id" = ?`, m.ID, m.AccountID).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("teamRepo.Update: %w", err)
	}
	return nil
}

func (r *TeamRepo) Delete(ctx context.Context, accountID, id int64) error {
	_, err := r.DB().NewDelete().TableExpr(`"teams"`).
		Where(`"id" = ? AND "account_id" = ?`, id, accountID).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("teamRepo.Delete: %w", err)
	}
	return nil
}

func (r *TeamRepo) AddMember(ctx context.Context, teamID, userID int64) error {
	m := &models.TeamMember{TeamID: teamID, UserID: userID}
	_, err := r.DB().NewInsert().Model(m).On("CONFLICT DO NOTHING").Exec(ctx)
	return err
}

func (r *TeamRepo) RemoveMember(ctx context.Context, teamID, userID int64) error {
	_, err := r.DB().NewDelete().TableExpr(`"team_members"`).
		Where(`"team_id" = ? AND "user_id" = ?`, teamID, userID).
		Exec(ctx)
	return err
}
