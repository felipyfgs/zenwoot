package handler

import (
	"github.com/gofiber/fiber/v2"

	"wzap/internal/domain"
	"wzap/internal/dto"
	"wzap/internal/service"
)

type TeamHandler struct {
	teamSvc *service.TeamService
}

func NewTeamHandler(teamSvc *service.TeamService) *TeamHandler {
	return &TeamHandler{teamSvc: teamSvc}
}

// Create godoc
// @Summary     Create a new team
// @Description Creates a new team within the account
// @Tags        Teams
// @Accept      json
// @Produce     json
// @Param       body body     dto.TeamCreateReq true "Team data"
// @Success     200  {object} dto.APIResponse
// @Security    ApiKey
// @Router      /api/v1/teams [post]
func (h *TeamHandler) Create(c *fiber.Ctx) error {
	accountID := getAccountID(c)

	var req dto.TeamCreateReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResp("Bad Request", err.Error()))
	}

	team, err := h.teamSvc.Create(c.Context(), accountID, service.CreateTeamReq{
		Name:            req.Name,
		Description:     req.Description,
		AllowAutoAssign: req.AllowAutoAssign,
	})
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResp("Bad Request", err.Error()))
	}

	return c.Status(fiber.StatusCreated).JSON(dto.SuccessResp(team))
}

// List godoc
// @Summary     List teams
// @Description Returns all teams for the account
// @Tags        Teams
// @Produce     json
// @Success     200 {object} dto.APIResponse
// @Security    ApiKey
// @Router      /api/v1/teams [get]
func (h *TeamHandler) List(c *fiber.Ctx) error {
	accountID := getAccountID(c)

	teams, err := h.teamSvc.ListByAccount(c.Context(), accountID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResp("Internal Server Error", err.Error()))
	}
	return c.JSON(dto.SuccessResp(teams))
}

// Get godoc
// @Summary     Get team
// @Description Returns team by ID
// @Tags        Teams
// @Produce     json
// @Param       teamId path string true "Team ID"
// @Success     200 {object} dto.APIResponse
// @Security    ApiKey
// @Router      /api/v1/teams/{teamId} [get]
func (h *TeamHandler) Get(c *fiber.Ctx) error {
	teamID := c.Params("teamId")
	team, err := h.teamSvc.Get(c.Context(), teamID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResp("Not Found", err.Error()))
	}
	return c.JSON(dto.SuccessResp(team))
}

// Update godoc
// @Summary     Update team
// @Description Updates team data
// @Tags        Teams
// @Accept      json
// @Produce     json
// @Param       teamId path string true "Team ID"
// @Param       body   body  dto.TeamUpdateReq true "Team data"
// @Success     200 {object} dto.APIResponse
// @Security    ApiKey
// @Router      /api/v1/teams/{teamId} [put]
func (h *TeamHandler) Update(c *fiber.Ctx) error {
	teamID := c.Params("teamId")

	var req dto.TeamUpdateReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResp("Bad Request", err.Error()))
	}

	team, err := h.teamSvc.Update(c.Context(), teamID, service.UpdateTeamReq{
		Name:            req.Name,
		Description:     req.Description,
		AllowAutoAssign: req.AllowAutoAssign,
	})
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResp("Bad Request", err.Error()))
	}
	return c.JSON(dto.SuccessResp(team))
}

// Delete godoc
// @Summary     Delete team
// @Description Removes team from account
// @Tags        Teams
// @Produce     json
// @Param       teamId path string true "Team ID"
// @Success     200 {object} dto.APIResponse
// @Security    ApiKey
// @Router      /api/v1/teams/{teamId} [delete]
func (h *TeamHandler) Delete(c *fiber.Ctx) error {
	teamID := c.Params("teamId")
	if err := h.teamSvc.Delete(c.Context(), teamID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResp("Internal Server Error", err.Error()))
	}
	return c.JSON(dto.SuccessResp(nil))
}

// AddMember godoc
// @Summary     Add member to team
// @Description Adds a user to a team
// @Tags        Teams
// @Produce     json
// @Param       teamId path string true "Team ID"
// @Param       userId path string true "User ID"
// @Success     200 {object} dto.APIResponse
// @Security    ApiKey
// @Router      /api/v1/teams/{teamId}/members/{userId} [post]
func (h *TeamHandler) AddMember(c *fiber.Ctx) error {
	teamID := c.Params("teamId")
	userID := c.Params("userId")

	var req dto.TeamMemberReq
	_ = c.BodyParser(&req) // optional body

	role := domain.RoleAgent
	if req.Role != "" {
		role = domain.UserRole(req.Role)
	}

	if err := h.teamSvc.AddMember(c.Context(), teamID, userID, role); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResp("Bad Request", err.Error()))
	}
	return c.JSON(dto.SuccessResp(nil))
}

// RemoveMember godoc
// @Summary     Remove member from team
// @Description Removes a user from a team
// @Tags        Teams
// @Produce     json
// @Param       teamId path string true "Team ID"
// @Param       userId path string true "User ID"
// @Success     200 {object} dto.APIResponse
// @Security    ApiKey
// @Router      /api/v1/teams/{teamId}/members/{userId} [delete]
func (h *TeamHandler) RemoveMember(c *fiber.Ctx) error {
	teamID := c.Params("teamId")
	userID := c.Params("userId")

	if err := h.teamSvc.RemoveMember(c.Context(), teamID, userID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResp("Bad Request", err.Error()))
	}
	return c.JSON(dto.SuccessResp(nil))
}

// ListMembers godoc
// @Summary     List team members
// @Description Returns all members of a team
// @Tags        Teams
// @Produce     json
// @Param       teamId path string true "Team ID"
// @Success     200 {object} dto.APIResponse
// @Security    ApiKey
// @Router      /api/v1/teams/{teamId}/members [get]
func (h *TeamHandler) ListMembers(c *fiber.Ctx) error {
	teamID := c.Params("teamId")

	members, err := h.teamSvc.ListMembers(c.Context(), teamID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResp("Internal Server Error", err.Error()))
	}
	return c.JSON(dto.SuccessResp(members))
}
