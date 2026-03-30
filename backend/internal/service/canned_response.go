package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"wzap/internal/domain"
	"wzap/internal/repo"
)

type CannedResponseService struct {
	repo *repo.CannedResponseRepository
}

func NewCannedResponseService(repo *repo.CannedResponseRepository) *CannedResponseService {
	return &CannedResponseService{repo: repo}
}

type CreateCannedResponseReq struct {
	ShortCode string `json:"shortCode"`
	Content   string `json:"content"`
}

func (s *CannedResponseService) Create(ctx context.Context, accountID string, req CreateCannedResponseReq) (*domain.CannedResponse, error) {
	now := time.Now()
	cr := &domain.CannedResponse{
		ID:        uuid.NewString(),
		AccountID: accountID,
		ShortCode: req.ShortCode,
		Content:   req.Content,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.repo.Create(ctx, cr); err != nil {
		return nil, err
	}
	return cr, nil
}

func (s *CannedResponseService) Get(ctx context.Context, id string) (*domain.CannedResponse, error) {
	cr, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, notFoundErrorf("canned response not found")
	}
	return cr, nil
}

func (s *CannedResponseService) List(ctx context.Context, accountID string) ([]domain.CannedResponse, error) {
	return s.repo.FindByAccountID(ctx, accountID)
}

func (s *CannedResponseService) Search(ctx context.Context, accountID, query string) ([]domain.CannedResponse, error) {
	if query == "" {
		return s.repo.FindByAccountID(ctx, accountID)
	}
	return s.repo.Search(ctx, accountID, query)
}

type UpdateCannedResponseReq struct {
	ShortCode string `json:"shortCode,omitempty"`
	Content   string `json:"content,omitempty"`
}

func (s *CannedResponseService) Update(ctx context.Context, id string, req UpdateCannedResponseReq) (*domain.CannedResponse, error) {
	cr, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, notFoundErrorf("canned response not found")
	}

	if req.ShortCode != "" {
		cr.ShortCode = req.ShortCode
	}
	if req.Content != "" {
		cr.Content = req.Content
	}

	if err := s.repo.Update(ctx, cr); err != nil {
		return nil, err
	}
	return cr, nil
}

func (s *CannedResponseService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
