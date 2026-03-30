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

	for attempts := 0; attempts < 3; attempts++ {
		var members []models.InboxMember
		err := s.db.NewSelect().Model(&members).
			Where(`"inbox_id" = ?`, conv.InboxID).
			Scan(ctx)
		if err != nil {
			return nil, fmt.Errorf("assignmentService.AutoAssign: %w", err)
		}
		if len(members) == 0 {
			return nil, fmt.Errorf("assignmentService.AutoAssign: no agents available")
		}

		picked := members[rand.Intn(len(members))]

		_, err = s.db.NewUpdate().Model(conv).
			Set(`"assignee_id" = ?`, picked.UserID).
			Where(`"id" = ? AND "assignee_id" IS NULL`, conv.ID).
			Exec(ctx)
		if err != nil {
			return nil, fmt.Errorf("assignmentService.AutoAssign: %w", err)
		}

		conv.AssigneeID = &picked.UserID
		return &picked.UserID, nil
	}

	return nil, fmt.Errorf("assignmentService.AutoAssign: failed after 3 attempts due to race condition")
}
