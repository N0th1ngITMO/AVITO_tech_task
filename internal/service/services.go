package services

import (
	"context"
	models "pr_task/internal/model"
)

type Service interface {
	CreateTeam(ctx context.Context, team models.Team) (*models.Team, error)
	GetTeam(ctx context.Context, teamName string) (*models.Team, error)

	SetUserActive(ctx context.Context, userID string, isActive bool) (*models.User, error)
	GetUserReviewPRs(ctx context.Context, userID string) (*models.UserReviewResponse, error)
	MassDeactivateTeamUsers(ctx context.Context, teamName string, excludeUserIDs []string) (*MassDeactivationResponse, error)

	CreatePullRequest(ctx context.Context, prID, name, authorID string) (*models.PullRequest, error)
	MergePullRequest(ctx context.Context, prID string) (*models.PullRequest, error)
	ReassignReviewer(ctx context.Context, prID, oldUserID string) (*models.ReassignResponse, error)

	GetUserReviewStats(ctx context.Context) ([]UserReviewStatsResponse, error)
	GetPRReviewStats(ctx context.Context) ([]PRReviewStatsResponse, error)
	GetOverallStats(ctx context.Context) (*OverallStatsResponse, error)
}
