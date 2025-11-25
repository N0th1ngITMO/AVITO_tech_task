package handler

import (
	"net/http"
	"pr_task/internal/dto"
	"time"

	"github.com/labstack/echo/v4"
)

// HealthCheck проверка здоровья сервиса
// @Summary Проверка здоровья
// @Description Возвращает статус сервиса
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} HealthResponse "Статус сервиса"
// @Router /health [get]
func (h *Handler) HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, dto.HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().Format(time.RFC3339),
		Service:   "PR Reviewer Assignment Service",
	})
}
