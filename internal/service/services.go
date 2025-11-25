package services

import (
	"context"
	"pr_task/internal/dto"
	models "pr_task/internal/model"
)

type Service interface {
	CreateTeam(ctx context.Context, team models.Team) (*models.Team, error)
	GetTeam(ctx context.Context, teamName string) (*models.Team, error)

	SetUserActive(ctx context.Context, userID string, isActive bool) (*models.User, error)
	GetUserReviewPRs(ctx context.Context, userID string) (*dto.UserReviewResponse, error)
	MassDeactivateTeamUsers(ctx context.Context, teamName string, excludeUserIDs []string) (*dto.MassDeactivationResponse, error)

	CreatePullRequest(ctx context.Context, prID, name, authorID string) (*dto.PullRequest, error)
	MergePullRequest(ctx context.Context, prID string) (*dto.PullRequest, error)
	ReassignReviewer(ctx context.Context, prID, oldUserID string) (*dto.ReassignResponse, error)

	GetUserReviewStats(ctx context.Context) ([]dto.UserReviewStatsResponse, error)
	GetPRReviewStats(ctx context.Context) ([]dto.PRReviewStatsResponse, error)
	GetOverallStats(ctx context.Context) (*dto.OverallStatsResponse, error)
}
