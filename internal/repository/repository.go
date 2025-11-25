package repository

import (
	"context"
	"pr_task/internal/dto"
	"time"

	"pr_task/internal/model"
)

type Repository interface {
	CreateTeam(ctx context.Context, team models.Team) error
	GetTeam(ctx context.Context, teamName string) (*models.Team, error)
	TeamExists(ctx context.Context, teamName string) (bool, error)

	CreateOrUpdateUser(ctx context.Context, user models.TeamMember, teamName string) error
	GetUser(ctx context.Context, userID string) (*models.User, error)
	UpdateUserActive(ctx context.Context, userID string, isActive bool) error
	GetActiveTeamMembers(ctx context.Context, teamName string, excludeUserID string) ([]models.User, error)
	GetRandomActiveTeamMember(ctx context.Context, teamName string, excludeUserIDs []string) (*models.User, error)
	MassDeactivateUsers(ctx context.Context, teamName string, excludeUserIDs []string) (int, error)
	GetOpenPRsWithReviewers(ctx context.Context, teamName string) ([]models.OpenPRInfo, error)
	UpdatePRReviewersBatch(ctx context.Context, updates []models.PRReviewersUpdate) error

	CreatePR(ctx context.Context, pr dto.PullRequest) error
	GetPR(ctx context.Context, prID string) (*dto.PullRequest, error)
	PRExists(ctx context.Context, prID string) (bool, error)
	UpdatePRStatus(ctx context.Context, prID string, status string, mergedAt *time.Time) error
	UpdatePRReviewers(ctx context.Context, prID string, reviewers []string) error
	GetPRsByReviewer(ctx context.Context, userID string) ([]dto.PullRequest, error)

	GetReviewStatsByUser(ctx context.Context) ([]models.UserReviewStats, error)
	GetReviewStatsByPR(ctx context.Context) ([]models.PRReviewStats, error)
	GetOverallStats(ctx context.Context) (*models.OverallStats, error)
}
