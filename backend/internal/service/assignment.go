package service

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"

	"wzap/internal/domain"
	"wzap/internal/repo"
	"wzap/internal/ws"
)

// AssignmentService handles conversation assignment logic.
type AssignmentService struct {
	convRepo        *repo.ConversationRepository
	inboxMemberRepo *repo.InboxMemberRepository
	teamMemberRepo  *repo.TeamMemberRepository
	wsHub           *ws.Hub
	accountID       string
}

func NewAssignmentService(
	convRepo *repo.ConversationRepository,
	inboxMemberRepo *repo.InboxMemberRepository,
	teamMemberRepo *repo.TeamMemberRepository,
	wsHub *ws.Hub,
	accountID string,
) *AssignmentService {
	return &AssignmentService{
		convRepo:        convRepo,
		inboxMemberRepo: inboxMemberRepo,
		teamMemberRepo:  teamMemberRepo,
		wsHub:           wsHub,
		accountID:       accountID,
	}
}

// AssignToUser assigns a conversation to a specific user.
func (s *AssignmentService) AssignToUser(ctx context.Context, conversationID, userID string) (*domain.Conversation, error) {
	conv, err := s.convRepo.FindByID(ctx, conversationID)
	if err != nil {
		return nil, notFoundErrorf("conversation not found")
	}

	// Update assignee
	now := time.Now()
	if err := s.updateAssignee(ctx, conversationID, &userID, nil); err != nil {
		return nil, err
	}

	conv.AssigneeID = userID
	conv.UpdatedAt = now

	// Publish event
	if s.wsHub != nil {
		s.wsHub.Publish(s.accountID, conv.InboxID, "conversation.assigned", map[string]interface{}{
			"conversationId": conversationID,
			"assigneeId":     userID,
		})
	}

	log.Info().Str("conversationId", conversationID).Str("assigneeId", userID).Msg("Conversation assigned to user")
	return conv, nil
}

// AssignToTeam assigns a conversation to a team.
func (s *AssignmentService) AssignToTeam(ctx context.Context, conversationID, teamID string) (*domain.Conversation, error) {
	conv, err := s.convRepo.FindByID(ctx, conversationID)
	if err != nil {
		return nil, notFoundErrorf("conversation not found")
	}

	now := time.Now()
	if err := s.updateAssignee(ctx, conversationID, nil, &teamID); err != nil {
		return nil, err
	}

	conv.TeamID = teamID
	conv.UpdatedAt = now

	if s.wsHub != nil {
		s.wsHub.Publish(s.accountID, conv.InboxID, "conversation.team_assigned", map[string]interface{}{
			"conversationId": conversationID,
			"teamId":         teamID,
		})
	}

	log.Info().Str("conversationId", conversationID).Str("teamId", teamID).Msg("Conversation assigned to team")
	return conv, nil
}

// Unassign removes assignment from a conversation.
func (s *AssignmentService) Unassign(ctx context.Context, conversationID string) (*domain.Conversation, error) {
	conv, err := s.convRepo.FindByID(ctx, conversationID)
	if err != nil {
		return nil, notFoundErrorf("conversation not found")
	}

	if err := s.updateAssignee(ctx, conversationID, nil, nil); err != nil {
		return nil, err
	}

	conv.AssigneeID = ""
	conv.TeamID = ""

	if s.wsHub != nil {
		s.wsHub.Publish(s.accountID, conv.InboxID, "conversation.unassigned", map[string]interface{}{
			"conversationId": conversationID,
		})
	}

	return conv, nil
}

// AutoAssign performs automatic round-robin assignment based on inbox members.
// It picks the member with the oldest lastAssignedAt (or never assigned), then updates the timestamp.
func (s *AssignmentService) AutoAssign(ctx context.Context, conversationID string) (*domain.Conversation, error) {
	conv, err := s.convRepo.FindByID(ctx, conversationID)
	if err != nil {
		return nil, notFoundErrorf("conversation not found")
	}

	// Get next member via round-robin (least recently assigned)
	member, err := s.inboxMemberRepo.FindNextRoundRobin(ctx, conv.InboxID)
	if err != nil {
		log.Debug().Str("inboxId", conv.InboxID).Msg("No inbox members for auto-assignment")
		return conv, nil
	}

	// Update round-robin tracking
	if err := s.inboxMemberRepo.UpdateLastAssigned(ctx, conv.InboxID, member.UserID); err != nil {
		log.Warn().Err(err).Str("userId", member.UserID).Msg("Failed to update last assigned timestamp")
	}

	return s.AssignToUser(ctx, conversationID, member.UserID)
}

// AutoAssignToTeam performs automatic assignment to a team member.
func (s *AssignmentService) AutoAssignToTeam(ctx context.Context, conversationID, teamID string) (*domain.Conversation, error) {
	// Get team members
	members, err := s.teamMemberRepo.FindByTeamID(ctx, teamID)
	if err != nil || len(members) == 0 {
		log.Debug().Str("teamId", teamID).Msg("No team members for auto-assignment")
		return nil, nil
	}

	// Simple round-robin: pick first member
	assigneeID := members[0].UserID

	// Set both team and assignee
	conv, err := s.convRepo.FindByID(ctx, conversationID)
	if err != nil {
		return nil, notFoundErrorf("conversation not found")
	}

	now := time.Now()
	if err := s.updateAssignee(ctx, conversationID, &assigneeID, &teamID); err != nil {
		return nil, err
	}

	conv.AssigneeID = assigneeID
	conv.TeamID = teamID
	conv.UpdatedAt = now

	return conv, nil
}

func (s *AssignmentService) updateAssignee(ctx context.Context, conversationID string, userID, teamID *string) error {
	return s.convRepo.UpdateAssignee(ctx, conversationID, userID, teamID)
}

// UpdatePriority updates the priority of a conversation.
func (s *AssignmentService) UpdatePriority(ctx context.Context, conversationID string, priority domain.ConversationPriority) (*domain.Conversation, error) {
	conv, err := s.convRepo.FindByID(ctx, conversationID)
	if err != nil {
		return nil, notFoundErrorf("conversation not found")
	}

	if err := s.convRepo.UpdatePriority(ctx, conversationID, priority); err != nil {
		return nil, err
	}

	conv.Priority = priority

	if s.wsHub != nil {
		s.wsHub.Publish(s.accountID, conv.InboxID, "conversation.priority_updated", map[string]interface{}{
			"conversationId": conversationID,
			"priority":       string(priority),
		})
	}

	return conv, nil
}

// Snooze snoozes a conversation until a specific time.
func (s *AssignmentService) Snooze(ctx context.Context, conversationID string, until time.Time) (*domain.Conversation, error) {
	conv, err := s.convRepo.FindByID(ctx, conversationID)
	if err != nil {
		return nil, notFoundErrorf("conversation not found")
	}

	if err := s.convRepo.Snooze(ctx, conversationID, until); err != nil {
		return nil, err
	}

	conv.Status = domain.ConversationSnoozed
	conv.SnoozedUntil = &until

	if s.wsHub != nil {
		s.wsHub.Publish(s.accountID, conv.InboxID, "conversation.snoozed", map[string]interface{}{
			"conversationId": conversationID,
			"snoozedUntil":   until.Format(time.RFC3339),
		})
	}

	return conv, nil
}

// UnSnooze unsnoozes a conversation.
func (s *AssignmentService) UnSnooze(ctx context.Context, conversationID string) (*domain.Conversation, error) {
	conv, err := s.convRepo.FindByID(ctx, conversationID)
	if err != nil {
		return nil, notFoundErrorf("conversation not found")
	}

	if err := s.convRepo.UnSnooze(ctx, conversationID); err != nil {
		return nil, err
	}

	conv.Status = domain.ConversationOpen
	conv.SnoozedUntil = nil

	return conv, nil
}
