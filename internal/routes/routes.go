package routes

import (
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	handlers "pr_task/internal/handler"
)

func RegisterRoutes(e *echo.Echo, handler *handlers.Handler) {
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.POST("/team/add", handler.AddTeam)
	e.GET("/team/get", handler.GetTeam)

	e.POST("/users/setIsActive", handler.SetUserActive)
	e.GET("/users/getReview", handler.GetUserReviews)
	e.POST("/users/massDeactivate", handler.MassDeactivateTeamUsers)

	e.POST("/pullRequest/create", handler.CreatePR)
	e.POST("/pullRequest/merge", handler.MergePR)
	e.POST("/pullRequest/reassign", handler.ReassignReviewer)

	e.GET("/stats/users", handler.GetUserReviewStats)
	e.GET("/stats/prs", handler.GetPRReviewStats)
	e.GET("/stats/overall", handler.GetOverallStats)

	e.GET("/health", handler.HealthCheck)
}
