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
// @Success 200 {object} UserReviewStatsResponse "Статистика по пользователям"
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
// @Success 200 {object} PRReviewStatsResponse "Статистика по PR"
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
// @Success 200 {object} OverallStatsResponse "Общая статистика"
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

type UserReviewStatsResponse struct {
	Stats []struct {
		UserID      string `json:"user_id" example:"u1"`
		Username    string `json:"username" example:"Alice"`
		TeamName    string `json:"team_name" example:"backend"`
		ReviewCount int    `json:"review_count" example:"5"`
	} `json:"stats"`
}

type PRReviewStatsResponse struct {
	Stats []struct {
		PRID          string `json:"pr_id" example:"pr-1001"`
		PRName        string `json:"pr_name" example:"Add search feature"`
		AuthorID      string `json:"author_id" example:"u1"`
		Status        string `json:"status" example:"OPEN"`
		ReviewerCount int    `json:"reviewer_count" example:"2"`
	} `json:"stats"`
}

type OverallStatsResponse struct {
	Stats struct {
		TotalPRs        int     `json:"total_prs" example:"10"`
		OpenPRs         int     `json:"open_prs" example:"7"`
		MergedPRs       int     `json:"merged_prs" example:"3"`
		TotalUsers      int     `json:"total_users" example:"15"`
		ActiveUsers     int     `json:"active_users" example:"12"`
		TotalReviews    int     `json:"total_reviews" example:"20"`
		AvgReviewsPerPR float64 `json:"avg_reviews_per_pr" example:"2.0"`
	} `json:"stats"`
}
