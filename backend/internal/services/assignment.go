package services

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/uptrace/bun"

	"github.com/felipyfgs/zenwoot/backend/internal/models"
)

type AssignmentService struct {
	db *bun.DB
}

func NewAssignmentService(db *bun.DB) *AssignmentService {
	return &AssignmentService{db: db}
}

func (s *AssignmentService) AutoAssign(ctx context.Context, conv *models.Conversation) (*int64, error) {
	if conv.AssigneeID != nil {
		return conv.AssigneeID, nil
	}
	var members []models.InboxMember
	err := s.db.NewSelect().Model(&members).
		Where(`"inbox_id" = ?`, conv.InboxID).
		Scan(ctx)
	if err != nil || len(members) == 0 {
		return nil, fmt.Errorf("assignmentService.AutoAssign: no agents available")
	}
	picked := members[rand.Intn(len(members))]
	return &picked.UserID, nil
}
