package service

import (
	"context"

	"wzap/internal/domain"
	"wzap/internal/dto"
	"wzap/internal/repo"
)

// ContactCatalogService handles the persisted contacts catalog (wzContacts table).
// This is separate from ContactService which operates on live WhatsApp state.
type ContactCatalogService struct {
	contactRepo *repo.ContactRepository
	convRepo    *repo.ConversationRepository
}

func NewContactCatalogService(contactRepo *repo.ContactRepository, convRepo *repo.ConversationRepository) *ContactCatalogService {
	return &ContactCatalogService{contactRepo: contactRepo, convRepo: convRepo}
}

func (s *ContactCatalogService) List(ctx context.Context, accountID string, req dto.PaginationReq) ([]domain.Contact, error) {
	limit, offset := req.LimitOffset()
	return s.contactRepo.FindByAccountID(ctx, accountID, limit, offset)
}

func (s *ContactCatalogService) Get(ctx context.Context, id string) (*domain.Contact, error) {
	c, err := s.contactRepo.FindByID(ctx, id)
	if err != nil {
		return nil, notFoundErrorf("contact not found")
	}
	return c, nil
}

func (s *ContactCatalogService) ListConversations(ctx context.Context, contactID string, req dto.PaginationReq) ([]domain.Conversation, error) {
	limit, offset := req.LimitOffset()
	return s.convRepo.FindByContactID(ctx, contactID, limit, offset)
}
