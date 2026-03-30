package services

import (
	"context"
	"encoding/json"

	"github.com/nats-io/nats.go"

	"github.com/felipyfgs/zenwoot/backend/internal/models"
	"github.com/felipyfgs/zenwoot/backend/internal/repo"
)

type ContactService struct {
	contactRepo *repo.ContactRepo
	nc          *nats.Conn
}

func NewContactService(contactRepo *repo.ContactRepo, nc *nats.Conn) *ContactService {
	return &ContactService{contactRepo: contactRepo, nc: nc}
}

func (s *ContactService) GetByID(ctx context.Context, accountID, id int64) (*models.Contact, error) {
	return s.contactRepo.GetByID(ctx, accountID, id)
}

func (s *ContactService) Search(ctx context.Context, accountID int64, q string, page, pageSize int) ([]*models.Contact, int, error) {
	return s.contactRepo.Search(ctx, accountID, q, page, pageSize)
}

func (s *ContactService) Create(ctx context.Context, m *models.Contact) (*models.Contact, error) {
	if err := s.contactRepo.Create(ctx, m); err != nil {
		return nil, err
	}
	s.publish("zenwoot.contact.created", m)
	return m, nil
}

func (s *ContactService) Update(ctx context.Context, m *models.Contact) (*models.Contact, error) {
	if err := s.contactRepo.Update(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *ContactService) publish(subject string, payload any) {
	if s.nc == nil {
		return
	}
	data, _ := json.Marshal(payload)
	_ = s.nc.Publish(subject, data)
}
