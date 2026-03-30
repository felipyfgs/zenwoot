package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"wzap/internal/domain"
	"wzap/internal/repo"
)

type UserService struct {
	userRepo        *repo.UserRepository
	accountUserRepo *repo.AccountUserRepository
}

func NewUserService(userRepo *repo.UserRepository, accountUserRepo *repo.AccountUserRepository) *UserService {
	return &UserService{userRepo: userRepo, accountUserRepo: accountUserRepo}
}

type CreateUserReq struct {
	Email       string          `json:"email"`
	Name        string          `json:"name"`
	DisplayName string          `json:"displayName,omitempty"`
	Role        domain.UserRole `json:"role"`
	Password    string          `json:"password,omitempty"`
	Provider    string          `json:"provider,omitempty"`
	UID         string          `json:"uid,omitempty"`
}

func (s *UserService) Create(ctx context.Context, accountID string, req CreateUserReq) (*domain.User, error) {
	now := time.Now()
	user := &domain.User{
		ID:          uuid.NewString(),
		Email:       req.Email,
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Role:        req.Role,
		Status:      domain.UserActive,
		Provider:    req.Provider,
		UID:         req.UID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if req.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		user.PasswordHash = string(hash)
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Create account-user link
	accountUser := &domain.AccountUser{
		ID:        uuid.NewString(),
		AccountID: accountID,
		UserID:    user.ID,
		Role:      req.Role,
		Active:    true,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.accountUserRepo.Create(ctx, accountUser); err != nil {
		return nil, fmt.Errorf("failed to link user to account: %w", err)
	}

	return user, nil
}

func (s *UserService) Get(ctx context.Context, id string) (*domain.User, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, notFoundErrorf("user not found")
	}
	return user, nil
}

func (s *UserService) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, notFoundErrorf("user not found")
	}
	return user, nil
}

func (s *UserService) ListByAccount(ctx context.Context, accountID string) ([]domain.User, error) {
	return s.userRepo.FindByAccountID(ctx, accountID)
}

type UpdateUserReq struct {
	Name        string            `json:"name,omitempty"`
	DisplayName string            `json:"displayName,omitempty"`
	AvatarURL   string            `json:"avatarUrl,omitempty"`
	Role        domain.UserRole   `json:"role,omitempty"`
	Status      domain.UserStatus `json:"status,omitempty"`
}

func (s *UserService) Update(ctx context.Context, id string, req UpdateUserReq) (*domain.User, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, notFoundErrorf("user not found")
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.DisplayName != "" {
		user.DisplayName = req.DisplayName
	}
	if req.AvatarURL != "" {
		user.AvatarURL = req.AvatarURL
	}
	if req.Role != "" {
		user.Role = req.Role
	}
	if req.Status != "" {
		user.Status = req.Status
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) Delete(ctx context.Context, id string) error {
	return s.userRepo.Delete(ctx, id)
}

// AccountUser operations

func (s *UserService) AddToAccount(ctx context.Context, accountID, userID string, role domain.UserRole) error {
	now := time.Now()
	accountUser := &domain.AccountUser{
		ID:        uuid.NewString(),
		AccountID: accountID,
		UserID:    userID,
		Role:      role,
		Active:    true,
		CreatedAt: now,
		UpdatedAt: now,
	}
	return s.accountUserRepo.Create(ctx, accountUser)
}

func (s *UserService) RemoveFromAccount(ctx context.Context, accountID, userID string) error {
	au, err := s.accountUserRepo.FindByAccountAndUser(ctx, accountID, userID)
	if err != nil {
		return err
	}
	return s.accountUserRepo.Delete(ctx, au.ID)
}

func (s *UserService) UpdateAccountRole(ctx context.Context, accountID, userID string, role domain.UserRole) error {
	au, err := s.accountUserRepo.FindByAccountAndUser(ctx, accountID, userID)
	if err != nil {
		return err
	}
	return s.accountUserRepo.UpdateRole(ctx, au.ID, role)
}
