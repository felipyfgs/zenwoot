package repo

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"wzap/internal/domain"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

const userCols = `"id","email","name",COALESCE("displayName",''),COALESCE("avatarUrl",''),
	"role","status","provider",COALESCE("uid",''),"settings","createdAt","updatedAt"`

func scanUser(row interface{ Scan(...any) error }) (*domain.User, error) {
	var u domain.User
	err := row.Scan(
		&u.ID, &u.Email, &u.Name, &u.DisplayName, &u.AvatarURL,
		&u.Role, &u.Status, &u.Provider, &u.UID, &u.Settings,
		&u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) Create(ctx context.Context, u *domain.User) error {
	q := `INSERT INTO "wzUsers" ("id","email","name","displayName","avatarUrl","role","status","passwordHash","provider","uid","settings","createdAt","updatedAt")
		  VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`
	_, err := r.db.Exec(ctx, q,
		u.ID, u.Email, u.Name, u.DisplayName, u.AvatarURL,
		u.Role, u.Status, u.PasswordHash, u.Provider, u.UID, u.Settings,
		u.CreatedAt, u.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *UserRepository) CountAll(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM "wzUsers"`).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}
	return count, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	q := `SELECT ` + userCols + ` FROM "wzUsers" WHERE "id"=$1`
	u, err := scanUser(r.db.QueryRow(ctx, q, id))
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return u, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	q := `SELECT ` + userCols + ` FROM "wzUsers" WHERE "email"=$1`
	u, err := scanUser(r.db.QueryRow(ctx, q, email))
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return u, nil
}

func (r *UserRepository) FindByEmailWithHash(ctx context.Context, email string) (*domain.User, error) {
	q := `SELECT ` + userCols + `,"passwordHash" FROM "wzUsers" WHERE "email"=$1`
	row := r.db.QueryRow(ctx, q, email)
	var u domain.User
	err := row.Scan(
		&u.ID, &u.Email, &u.Name, &u.DisplayName, &u.AvatarURL,
		&u.Role, &u.Status, &u.Provider, &u.UID, &u.Settings,
		&u.CreatedAt, &u.UpdatedAt, &u.PasswordHash,
	)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return &u, nil
}

const userColsQualified = `u."id",u."email",u."name",COALESCE(u."displayName",''),COALESCE(u."avatarUrl",''),` +
	`u."role",u."status",u."provider",COALESCE(u."uid",''),u."settings",u."createdAt",u."updatedAt"`

func (r *UserRepository) FindByAccountID(ctx context.Context, accountID string) ([]domain.User, error) {
	q := `SELECT ` + userColsQualified + ` FROM "wzUsers" u
		  INNER JOIN "wzAccountUsers" au ON u."id" = au."userId"
		  WHERE au."accountId"=$1 AND au."active"=true
		  ORDER BY u."name" ASC`
	rows, err := r.db.Query(ctx, q, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var list []domain.User
	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, *u)
	}
	return list, rows.Err()
}

func (r *UserRepository) Update(ctx context.Context, u *domain.User) error {
	q := `UPDATE "wzUsers" SET "name"=$1,"displayName"=$2,"avatarUrl"=$3,"role"=$4,"status"=$5,"settings"=$6,"updatedAt"=NOW()
		  WHERE "id"=$7`
	_, err := r.db.Exec(ctx, q, u.Name, u.DisplayName, u.AvatarURL, u.Role, u.Status, u.Settings, u.ID)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM "wzUsers" WHERE "id"=$1`, id)
	return err
}

// AccountUser Repository methods

type AccountUserRepository struct {
	db *pgxpool.Pool
}

func NewAccountUserRepository(db *pgxpool.Pool) *AccountUserRepository {
	return &AccountUserRepository{db: db}
}

func (r *AccountUserRepository) Create(ctx context.Context, au *domain.AccountUser) error {
	q := `INSERT INTO "wzAccountUsers" ("id","accountId","userId","role","active","createdAt","updatedAt")
		  VALUES ($1,$2,$3,$4,$5,$6,$7)`
	_, err := r.db.Exec(ctx, q, au.ID, au.AccountID, au.UserID, au.Role, au.Active, au.CreatedAt, au.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create account user: %w", err)
	}
	return nil
}

func (r *AccountUserRepository) FindByAccountAndUser(ctx context.Context, accountID, userID string) (*domain.AccountUser, error) {
	q := `SELECT "id","accountId","userId","role","active","createdAt","updatedAt"
		  FROM "wzAccountUsers" WHERE "accountId"=$1 AND "userId"=$2`
	var au domain.AccountUser
	err := r.db.QueryRow(ctx, q, accountID, userID).Scan(
		&au.ID, &au.AccountID, &au.UserID, &au.Role, &au.Active, &au.CreatedAt, &au.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("account user not found: %w", err)
	}
	return &au, nil
}

func (r *AccountUserRepository) FindByAccountID(ctx context.Context, accountID string) ([]domain.AccountUser, error) {
	q := `SELECT "id","accountId","userId","role","active","createdAt","updatedAt"
		  FROM "wzAccountUsers" WHERE "accountId"=$1`
	rows, err := r.db.Query(ctx, q, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to query account users: %w", err)
	}
	defer rows.Close()

	var list []domain.AccountUser
	for rows.Next() {
		var au domain.AccountUser
		if err := rows.Scan(&au.ID, &au.AccountID, &au.UserID, &au.Role, &au.Active, &au.CreatedAt, &au.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, au)
	}
	return list, rows.Err()
}

func (r *AccountUserRepository) UpdateRole(ctx context.Context, id string, role domain.UserRole) error {
	_, err := r.db.Exec(ctx, `UPDATE "wzAccountUsers" SET "role"=$1,"updatedAt"=NOW() WHERE "id"=$2`, role, id)
	return err
}

func (r *AccountUserRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM "wzAccountUsers" WHERE "id"=$1`, id)
	return err
}
