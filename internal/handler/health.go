package handler

import (
	"net/http"
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
	return c.JSON(http.StatusOK, HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().Format(time.RFC3339),
		Service:   "PR Reviewer Assignment Service",
	})
}

type HealthResponse struct {
	Status    string `json:"status" example:"healthy"`
	Timestamp string `json:"timestamp" example:"2025-11-25T16:30:45Z"`
	Service   string `json:"service" example:"PR Reviewer Assignment Service"`
}
