package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"pr_task/internal/dto"
	errors "pr_task/internal/error"
	"pr_task/internal/model"
)

// AddTeam создает команду с участниками
// @Summary Создать команду с участниками
// @Description Создает команду и обновляет/создает пользователей
// @Tags Teams
// @Accept json
// @Produce json
// @Param request body models.Team true "Данные команды"
// @Success 201 {object} map[string]interface{} "Команда создана"
// @Failure 400 {object} errors.ErrorResponse "Команда уже существует"
// @Failure 500 {object} errors.ErrorResponse "Внутренняя ошибка сервера"
// @Router /team/add [post]
func (h *Handler) AddTeam(c echo.Context) error {
	var req dto.TeamCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errors.NewErrorResponse("INVALID_REQUEST", "Invalid request body"))
	}

	if req.TeamName == "" || len(req.Members) == 0 {
		return c.JSON(http.StatusBadRequest, errors.NewErrorResponse("INVALID_REQUEST", "Team name and members are required"))
	}

	team, err := h.Service.CreateTeam(c.Request().Context(), models.Team{TeamName: req.TeamName, Members: req.Members})
	if err != nil {
		if errors.Is(err, errors.ErrTeamExists) {
			return c.JSON(http.StatusBadRequest, errors.NewErrorResponse(errors.CodeTeamExists, "team_name already exists"))
		}
		return c.JSON(http.StatusInternalServerError, errors.NewErrorResponse("INTERNAL_ERROR", "Failed to create team"))
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"team": team,
	})
}

// GetTeam получает команду с участниками
// @Summary Получить команду с участниками
// @Description Возвращает команду по имени
// @Tags Teams
// @Accept json
// @Produce json
// @Param team_name query string true "Уникальное имя команды" example:"backend"
// @Success 200 {object} models.Team "Объект команды"
// @Failure 400 {object} errors.ErrorResponse "Неверный запрос"
// @Failure 404 {object} errors.ErrorResponse "Команда не найдена"
// @Failure 500 {object} errors.ErrorResponse "Внутренняя ошибка сервера"
// @Router /team/get [get]
func (h *Handler) GetTeam(c echo.Context) error {
	teamName := c.QueryParam("team_name")
	if teamName == "" {
		return c.JSON(http.StatusBadRequest, errors.NewErrorResponse("INVALID_REQUEST", "team_name is required"))
	}

	team, err := h.Service.GetTeam(c.Request().Context(), teamName)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return c.JSON(http.StatusNotFound, errors.NewErrorResponse(errors.CodeNotFound, "Team not found"))
		}
		return c.JSON(http.StatusInternalServerError, errors.NewErrorResponse("INTERNAL_ERROR", "Failed to get team"))
	}

	return c.JSON(http.StatusOK, team)
}
