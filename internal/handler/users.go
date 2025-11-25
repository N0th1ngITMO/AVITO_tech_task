package handler

import (
	"net/http"
	"pr_task/internal/dto"
	errors "pr_task/internal/error"

	"github.com/labstack/echo/v4"
)

// SetUserActive устанавливает флаг активности пользователя
// @Summary Установить флаг активности пользователя
// @Description Обновляет активность пользователя
// @Tags Users
// @Accept json
// @Produce json
// @Param request body dto.SetUserActiveRequest true "Данные пользователя"
// @Success 200 {object} map[string]interface{} "Обновлённый пользователь"
// @Failure 400 {object} errors.ErrorResponse "Неверный запрос"
// @Failure 404 {object} errors.ErrorResponse "Пользователь не найден"
// @Failure 500 {object} errors.ErrorResponse "Внутренняя ошибка сервера"
// @Router /users/setIsActive [post]
func (h *Handler) SetUserActive(c echo.Context) error {
	var req dto.SetUserActiveRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errors.NewErrorResponse("INVALID_REQUEST", "Invalid request body"))
	}

	if req.UserID == "" {
		return c.JSON(http.StatusBadRequest, errors.NewErrorResponse("INVALID_REQUEST", "user_id is required"))
	}

	user, err := h.Service.SetUserActive(c.Request().Context(), req.UserID, req.IsActive)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return c.JSON(http.StatusNotFound, errors.NewErrorResponse(errors.CodeNotFound, "User not found"))
		}
		return c.JSON(http.StatusInternalServerError, errors.NewErrorResponse("INTERNAL_ERROR", "Failed to update user"))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"user": user,
	})
}

// GetUserReviews получает PR'ы, где пользователь назначен ревьювером
// @Summary Получить PR'ы, где пользователь назначен ревьювером
// @Description Возвращает список PR для ревью пользователя
// @Tags Users
// @Accept json
// @Produce json
// @Param user_id query string true "Идентификатор пользователя" example:"u1"
// @Success 200 {object} models.UserReviewResponse "Список PR'ов пользователя"
// @Failure 400 {object} errors.ErrorResponse "Неверный запрос"
// @Failure 404 {object} errors.ErrorResponse "Пользователь не найден"
// @Failure 500 {object} errors.ErrorResponse "Внутренняя ошибка сервера"
// @Router /users/getReview [get]
func (h *Handler) GetUserReviews(c echo.Context) error {
	userID := c.QueryParam("user_id")
	if userID == "" {
		return c.JSON(http.StatusBadRequest, errors.NewErrorResponse("INVALID_REQUEST", "user_id is required"))
	}

	result, err := h.Service.GetUserReviewPRs(c.Request().Context(), userID)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return c.JSON(http.StatusNotFound, errors.NewErrorResponse(errors.CodeNotFound, "User not found"))
		}
		return c.JSON(http.StatusInternalServerError, errors.NewErrorResponse("INTERNAL_ERROR", "Failed to get user reviews"))
	}

	return c.JSON(http.StatusOK, result)
}
