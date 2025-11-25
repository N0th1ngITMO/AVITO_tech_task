package handler

import (
	"fmt"
	"net/http"
	errors "pr_task/internal/error"

	"github.com/labstack/echo/v4"
)

// CreatePR @Summary Создать PR и автоматически назначить до 2 ревьюверов
// @Tags PullRequests
// @Accept json
// @Produce json
// @Param request body dto.CreatePullRequestRequest true "Данные PR"
// @Success 201 {object} map[string]interface{} "PR создан"
// @Failure 400 {object} errors.ErrorResponse "Неверный запрос"
// @Failure 404 {object} errors.ErrorResponse "Автор/команда не найдены"
// @Failure 409 {object} errors.ErrorResponse "PR уже существует"
// @Failure 500 {object} errors.ErrorResponse "Внутренняя ошибка сервера"
// @Router /pullRequest/create [post]
func (h *Handler) CreatePR(c echo.Context) error {
	var req struct {
		PullRequestID   string `json:"pull_request_id"`
		PullRequestName string `json:"pull_request_name"`
		AuthorID        string `json:"author_id"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errors.NewErrorResponse("INVALID_REQUEST", "Invalid request body"))
	}

	if req.PullRequestID == "" || req.PullRequestName == "" || req.AuthorID == "" {
		return c.JSON(http.StatusBadRequest, errors.NewErrorResponse("INVALID_REQUEST", "All fields are required"))
	}

	pr, err := h.Service.CreatePullRequest(c.Request().Context(), req.PullRequestID, req.PullRequestName, req.AuthorID)
	if err != nil {
		switch {
		case errors.Is(err, errors.ErrNotFound):
			return c.JSON(http.StatusNotFound, errors.NewErrorResponse(errors.CodeNotFound, "Author or team not found"))
		case errors.Is(err, errors.ErrPRExists):
			return c.JSON(http.StatusConflict, errors.NewErrorResponse(errors.CodePRExists, "PR id already exists"))
		default:
			return c.JSON(http.StatusInternalServerError, errors.NewErrorResponse("INTERNAL_ERROR", "Failed to create PR"))
		}
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"pr": pr,
	})
}

// MergePR помечает PR как MERGED
// @Summary Пометить PR как MERGED
// @Description Идемпотентная операция мержа PR
// @Tags PullRequests
// @Accept json
// @Produce json
// @Param request body dto.MergePullRequestRequest true "Данные PR"
// @Success 200 {object} models.PullRequest "PR в состоянии MERGED"
// @Failure 400 {object} errors.ErrorResponse "Неверный запрос"
// @Failure 404 {object} errors.ErrorResponse "PR не найден"
// @Failure 500 {object} errors.ErrorResponse "Внутренняя ошибка сервера"
// @Router /pullRequest/merge [post]
func (h *Handler) MergePR(c echo.Context) error {
	var req struct {
		PullRequestID string `json:"pull_request_id"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errors.NewErrorResponse("INVALID_REQUEST", "Invalid request body"))
	}

	if req.PullRequestID == "" {
		return c.JSON(http.StatusBadRequest, errors.NewErrorResponse("INVALID_REQUEST", "pull_request_id is required"))
	}

	pr, err := h.Service.MergePullRequest(c.Request().Context(), req.PullRequestID)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return c.JSON(http.StatusNotFound, errors.NewErrorResponse(errors.CodeNotFound, "PR not found"))
		}
		return c.JSON(http.StatusInternalServerError, errors.NewErrorResponse("INTERNAL_ERROR", "Failed to merge PR"))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"pr": pr,
	})
}

// ReassignReviewer переназначает ревьювера
// @Summary Переназначить конкретного ревьювера
// @Description Заменяет ревьювера на другого из его команды
// @Tags PullRequests
// @Accept json
// @Produce json
// @Param request body dto.ReassignReviewerRequest true "Данные для переназначения"
// @Success 200 {object} models.ReassignResponse "Переназначение выполнено"
// @Failure 400 {object} errors.ErrorResponse "Неверный запрос"
// @Failure 404 {object} errors.ErrorResponse "PR или пользователь не найден"
// @Failure 409 {object} errors.ErrorResponse "Нарушение доменных правил"
// @Failure 500 {object} errors.ErrorResponse "Внутренняя ошибка сервера"
// @Router /pullRequest/reassign [post]
func (h *Handler) ReassignReviewer(c echo.Context) error {
	var req struct {
		PullRequestID string `json:"pull_request_id"`
		OldReviewerID string `json:"old_reviewer_id"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errors.NewErrorResponse("INVALID_REQUEST", "Invalid request body"))
	}
	fmt.Printf("Bound req: %+v\n", req)
	if req.PullRequestID == "" || req.OldReviewerID == "" {
		return c.JSON(http.StatusBadRequest, errors.NewErrorResponse("INVALID_REQUEST", "Both pull_request_id and old_user_id are required"))
	}

	result, err := h.Service.ReassignReviewer(c.Request().Context(), req.PullRequestID, req.OldReviewerID)
	if err != nil {
		switch {
		case errors.Is(err, errors.ErrNotFound):
			return c.JSON(http.StatusNotFound, errors.NewErrorResponse(errors.CodeNotFound, "PR or user not found"))
		case errors.Is(err, errors.ErrPRMerged):
			return c.JSON(http.StatusConflict, errors.NewErrorResponse(errors.CodePRMerged, "cannot reassign on merged PR"))
		case errors.Is(err, errors.ErrNotAssigned):
			return c.JSON(http.StatusConflict, errors.NewErrorResponse(errors.CodeNotAssigned, "reviewer is not assigned to this PR"))
		case errors.Is(err, errors.ErrNoCandidate):
			return c.JSON(http.StatusConflict, errors.NewErrorResponse(errors.CodeNoCandidate, "no active replacement candidate in team"))
		default:
			return c.JSON(http.StatusInternalServerError, errors.NewErrorResponse("INTERNAL_ERROR", "Failed to reassign reviewer"))
		}
	}

	return c.JSON(http.StatusOK, result)
}
