package handler

import (
	"net/http"
	errors "pr_task/internal/error"

	"github.com/labstack/echo/v4"
)

// GetUserReviewStats возвращает статистику по ревьюверам
// @Summary Получить статистику по ревьюверам
// @Description Возвращает количество назначений по каждому пользователю
// @Tags Statistics
// @Accept json
// @Produce json
// @Success 200 {object} dto.UserReviewStatsResponse "Статистика по пользователям"
// @Failure 500 {object} errors.ErrorResponse "Внутренняя ошибка сервера"
// @Router /stats/users [get]
func (h *Handler) GetUserReviewStats(c echo.Context) error {
	stats, err := h.Service.GetUserReviewStats(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errors.NewErrorResponse("INTERNAL_ERROR", "Failed to get user review stats"))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"stats": stats,
	})
}

// GetPRReviewStats возвращает статистику по PR
// @Summary Получить статистику по PR
// @Description Возвращает количество ревьюверов по каждому PR
// @Tags Statistics
// @Accept json
// @Produce json
// @Success 200 {object} dto.PRReviewStatsResponse "Статистика по PR"
// @Failure 500 {object} errors.ErrorResponse "Внутренняя ошибка сервера"
// @Router /stats/prs [get]
func (h *Handler) GetPRReviewStats(c echo.Context) error {
	stats, err := h.Service.GetPRReviewStats(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errors.NewErrorResponse("INTERNAL_ERROR", "Failed to get PR review stats"))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"stats": stats,
	})
}

// GetOverallStats возвращает общую статистику
// @Summary Получить общую статистику
// @Description Возвращает общую статистику системы
// @Tags Statistics
// @Accept json
// @Produce json
// @Success 200 {object} dto.OverallStatsResponse "Общая статистика"
// @Failure 500 {object} errors.ErrorResponse "Внутренняя ошибка сервера"
// @Router /stats/overall [get]
func (h *Handler) GetOverallStats(c echo.Context) error {
	stats, err := h.Service.GetOverallStats(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errors.NewErrorResponse("INTERNAL_ERROR", "Failed to get overall stats"))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"stats": stats,
	})
}
