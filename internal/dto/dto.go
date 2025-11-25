package dto

type TeamCreateRequest struct {
	TeamName string       `json:"team_name" validate:"required"`
	Members  []TeamMember `json:"members" validate:"required,min=1"`
}

type TeamMember struct {
	UserID   string `json:"user_id" validate:"required"`
	Username string `json:"username" validate:"required"`
	IsActive bool   `json:"is_active"`
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
