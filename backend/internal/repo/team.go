package repo

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"wzap/internal/domain"
)

type TeamRepository struct {
	db *pgxpool.Pool
}

func NewTeamRepository(db *pgxpool.Pool) *TeamRepository {
	return &TeamRepository{db: db}
}

func (r *TeamRepository) Create(ctx context.Context, t *domain.Team) error {
	q := `INSERT INTO "wzTeams" ("id","accountId","name","description","allowAutoAssign","createdAt","updatedAt")
		  VALUES ($1,$2,$3,$4,$5,$6,$7)`
	_, err := r.db.Exec(ctx, q, t.ID, t.AccountID, t.Name, t.Description, t.AllowAutoAssign, t.CreatedAt, t.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create team: %w", err)
	}
	return nil
}

func (r *TeamRepository) FindByID(ctx context.Context, id string) (*domain.Team, error) {
	q := `SELECT "id","accountId","name",COALESCE("description",''),"allowAutoAssign","createdAt","updatedAt"
		  FROM "wzTeams" WHERE "id"=$1`
	var t domain.Team
	err := r.db.QueryRow(ctx, q, id).Scan(
		&t.ID, &t.AccountID, &t.Name, &t.Description, &t.AllowAutoAssign, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("team not found: %w", err)
	}
	return &t, nil
}

func (r *TeamRepository) FindByAccountID(ctx context.Context, accountID string) ([]domain.Team, error) {
	q := `SELECT "id","accountId","name",COALESCE("description",''),"allowAutoAssign","createdAt","updatedAt"
		  FROM "wzTeams" WHERE "accountId"=$1 ORDER BY "name" ASC`
	rows, err := r.db.Query(ctx, q, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to query teams: %w", err)
	}
	defer rows.Close()

	var list []domain.Team
	for rows.Next() {
		var t domain.Team
		if err := rows.Scan(&t.ID, &t.AccountID, &t.Name, &t.Description, &t.AllowAutoAssign, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, t)
	}
	return list, rows.Err()
}

func (r *TeamRepository) Update(ctx context.Context, t *domain.Team) error {
	q := `UPDATE "wzTeams" SET "name"=$1,"description"=$2,"allowAutoAssign"=$3,"updatedAt"=NOW() WHERE "id"=$4`
	_, err := r.db.Exec(ctx, q, t.Name, t.Description, t.AllowAutoAssign, t.ID)
	return err
}

func (r *TeamRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM "wzTeams" WHERE "id"=$1`, id)
	return err
}

// TeamMember Repository

type TeamMemberRepository struct {
	db *pgxpool.Pool
}

func NewTeamMemberRepository(db *pgxpool.Pool) *TeamMemberRepository {
	return &TeamMemberRepository{db: db}
}

func (r *TeamMemberRepository) Create(ctx context.Context, tm *domain.TeamMember) error {
	q := `INSERT INTO "wzTeamMembers" ("id","teamId","userId","role","createdAt")
		  VALUES ($1,$2,$3,$4,$5)`
	_, err := r.db.Exec(ctx, q, tm.ID, tm.TeamID, tm.UserID, tm.Role, tm.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create team member: %w", err)
	}
	return nil
}

func (r *TeamMemberRepository) FindByTeamID(ctx context.Context, teamID string) ([]domain.TeamMember, error) {
	q := `SELECT "id","teamId","userId","role","createdAt" FROM "wzTeamMembers" WHERE "teamId"=$1`
	rows, err := r.db.Query(ctx, q, teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to query team members: %w", err)
	}
	defer rows.Close()

	var list []domain.TeamMember
	for rows.Next() {
		var tm domain.TeamMember
		if err := rows.Scan(&tm.ID, &tm.TeamID, &tm.UserID, &tm.Role, &tm.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, tm)
	}
	return list, rows.Err()
}

func (r *TeamMemberRepository) Delete(ctx context.Context, teamID, userID string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM "wzTeamMembers" WHERE "teamId"=$1 AND "userId"=$2`, teamID, userID)
	return err
}

// InboxMember Repository

type InboxMemberRepository struct {
	db *pgxpool.Pool
}

func NewInboxMemberRepository(db *pgxpool.Pool) *InboxMemberRepository {
	return &InboxMemberRepository{db: db}
}

func (r *InboxMemberRepository) Create(ctx context.Context, im *domain.InboxMember) error {
	q := `INSERT INTO "wzInboxMembers" ("id","inboxId","userId","createdAt")
		  VALUES ($1,$2,$3,$4)`
	_, err := r.db.Exec(ctx, q, im.ID, im.InboxID, im.UserID, im.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create inbox member: %w", err)
	}
	return nil
}

func (r *InboxMemberRepository) FindByInboxID(ctx context.Context, inboxID string) ([]domain.InboxMember, error) {
	q := `SELECT "id","inboxId","userId","lastAssignedAt","createdAt" FROM "wzInboxMembers" WHERE "inboxId"=$1`
	rows, err := r.db.Query(ctx, q, inboxID)
	if err != nil {
		return nil, fmt.Errorf("failed to query inbox members: %w", err)
	}
	defer rows.Close()

	var list []domain.InboxMember
	for rows.Next() {
		var im domain.InboxMember
		if err := rows.Scan(&im.ID, &im.InboxID, &im.UserID, &im.LastAssignedAt, &im.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, im)
	}
	return list, rows.Err()
}

// FindNextRoundRobin returns the inbox member who was assigned least recently (round-robin).
func (r *InboxMemberRepository) FindNextRoundRobin(ctx context.Context, inboxID string) (*domain.InboxMember, error) {
	q := `SELECT "id","inboxId","userId","lastAssignedAt","createdAt"
		  FROM "wzInboxMembers"
		  WHERE "inboxId"=$1
		  ORDER BY "lastAssignedAt" ASC NULLS FIRST, "createdAt" ASC
		  LIMIT 1`
	var im domain.InboxMember
	err := r.db.QueryRow(ctx, q, inboxID).Scan(&im.ID, &im.InboxID, &im.UserID, &im.LastAssignedAt, &im.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("no inbox members available: %w", err)
	}
	return &im, nil
}

// UpdateLastAssigned updates the lastAssignedAt timestamp for a member.
func (r *InboxMemberRepository) UpdateLastAssigned(ctx context.Context, inboxID, userID string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE "wzInboxMembers" SET "lastAssignedAt"=NOW() WHERE "inboxId"=$1 AND "userId"=$2`,
		inboxID, userID,
	)
	return err
}

func (r *InboxMemberRepository) FindByUserID(ctx context.Context, userID string) ([]domain.InboxMember, error) {
	q := `SELECT "id","inboxId","userId","lastAssignedAt","createdAt" FROM "wzInboxMembers" WHERE "userId"=$1`
	rows, err := r.db.Query(ctx, q, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query inbox members: %w", err)
	}
	defer rows.Close()

	var list []domain.InboxMember
	for rows.Next() {
		var im domain.InboxMember
		if err := rows.Scan(&im.ID, &im.InboxID, &im.UserID, &im.LastAssignedAt, &im.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, im)
	}
	return list, rows.Err()
}

func (r *InboxMemberRepository) Delete(ctx context.Context, inboxID, userID string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM "wzInboxMembers" WHERE "inboxId"=$1 AND "userId"=$2`, inboxID, userID)
	return err
}
