package dto

import (
	models "pr_task/internal/model"
	"time"
)

type TeamCreateRequest struct {
	TeamName string              `json:"team_name" validate:"required"`
	Members  []models.TeamMember `json:"members" validate:"required,min=1"`
}

type SetUserActiveRequest struct {
	UserID   string `json:"user_id" validate:"required"`
	IsActive bool   `json:"is_active"`
}

type CreatePullRequestRequest struct {
	PullRequestID   string `json:"pull_request_id" validate:"required"`
	PullRequestName string `json:"pull_request_name" validate:"required"`
	AuthorID        string `json:"author_id" validate:"required"`
}

type MergePullRequestRequest struct {
	PullRequestID string `json:"pull_request_id" validate:"required"`
}

type ReassignReviewerRequest struct {
	PullRequestID string `json:"pull_request_id" validate:"required"`
	OldReviewerID string `json:"old_reviewer_id" validate:"required"`
}

type GetUserReviewRequest struct {
	UserID string `query:"user_id" validate:"required"`
}

type MassDeactivationRequest struct {
	TeamName       string   `json:"team_name" example:"backend"`
	ExcludeUserIDs []string `json:"exclude_user_ids,omitempty" example:"u1,u2"`
}

type MassDeactivationResponse struct {
	DeactivatedUsers int      `json:"deactivated_users" example:"5"`
	UpdatedPRs       int      `json:"updated_prs" example:"3"`
	FailedPRs        []string `json:"failed_prs,omitempty" example:"pr-1001"`
	ProcessingTime   int64    `json:"processing_time_ms" example:"85"`
}

type UserReviewStatsResponse struct {
	UserID      string `json:"user_id" example:"u1"`
	Username    string `json:"username" example:"Alice"`
	TeamName    string `json:"team_name" example:"backend"`
	ReviewCount int    `json:"review_count" example:"5"`
}

type PRReviewStatsResponse struct {
	PRID          string `json:"pr_id" example:"pr-1001"`
	PRName        string `json:"pr_name" example:"Add search feature"`
	AuthorID      string `json:"author_id" example:"u1"`
	Status        string `json:"status" example:"OPEN"`
	ReviewerCount int    `json:"reviewer_count" example:"2"`
}

type OverallStatsResponse struct {
	TotalPRs        int     `json:"total_prs" example:"10"`
	OpenPRs         int     `json:"open_prs" example:"7"`
	MergedPRs       int     `json:"merged_prs" example:"3"`
	TotalUsers      int     `json:"total_users" example:"15"`
	ActiveUsers     int     `json:"active_users" example:"12"`
	TotalReviews    int     `json:"total_reviews" example:"20"`
	AvgReviewsPerPR float64 `json:"avg_reviews_per_pr" example:"2.0"`
}

type HealthResponse struct {
	Status    string `json:"status" example:"healthy"`
	Timestamp string `json:"timestamp" example:"2025-11-25T16:30:45Z"`
	Service   string `json:"service" example:"PR Reviewer Assignment Service"`
}

type PullRequest struct {
	PullRequestID     string     `json:"pull_request_id"`
	PullRequestName   string     `json:"pull_request_name"`
	AuthorID          string     `json:"author_id"`
	Status            string     `json:"status"`
	AssignedReviewers []string   `json:"assigned_reviewers"`
	CreatedAt         *time.Time `json:"createdAt,omitempty"`
	MergedAt          *time.Time `json:"mergedAt,omitempty"`
}

type PullRequestShort struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
	Status          string `json:"status"`
}

type UserReviewResponse struct {
	UserID       string             `json:"user_id"`
	PullRequests []PullRequestShort `json:"pull_requests"`
}

type ReassignResponse struct {
	PR         *PullRequest `json:"pr"`
	ReplacedBy string       `json:"replaced_by"`
}
