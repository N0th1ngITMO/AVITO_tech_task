package handler

import (
	"net/http"
	errors "pr_task/internal/error"

	"github.com/labstack/echo/v4"
)

// MassDeactivateTeamUsers массово деактивирует пользователей команды
// @Summary Массовая деактивация пользователей команды
// @Description Деактивирует всех пользователей команды и безопасно переназначает ревьюверов открытых PR
// @Tags Users
// @Accept json
// @Produce json
// @Param request body MassDeactivationRequest true "Данные для массовой деактивации"
// @Success 200 {object} MassDeactivationResponse "Результат массовой деактивации"
// @Failure 400 {object} errors.ErrorResponse "Неверный запрос"
// @Failure 404 {object} errors.ErrorResponse "Команда не найдена"
// @Failure 500 {object} errors.ErrorResponse "Внутренняя ошибка сервера"
// @Router /users/massDeactivate [post]
func (h *Handler) MassDeactivateTeamUsers(c echo.Context) error {
	var req MassDeactivationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errors.NewErrorResponse("INVALID_REQUEST", "Invalid request body"))
	}

	if req.TeamName == "" {
		return c.JSON(http.StatusBadRequest, errors.NewErrorResponse("INVALID_REQUEST", "Team name is required"))
	}

	result, err := h.Service.MassDeactivateTeamUsers(c.Request().Context(), req.TeamName, req.ExcludeUserIDs)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errors.NewErrorResponse("INTERNAL_ERROR", err.Error()))
	}

	return c.JSON(http.StatusOK, result)
}

type MassDeactivationRequest struct {
	TeamName       string   `json:"team_name" example:"backend"`
	ExcludeUserIDs []string `json:"exclude_user_ids,omitempty" example:"u1,u2"`
}

type MassDeactivationResponse struct {
	DeactivatedUsers int      `json:"deactivated_users" example:"5"`
	UpdatedPRs       int      `json:"updated_prs" example:"3"`
	FailedPRs        []string `json:"failed_prs,omitempty" example:"pr-1001"`
	ProcessingTime   int64    `json:"processing_time_ms" example:"85"`
}
