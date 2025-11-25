package models

type Team struct {
	TeamName string       `json:"team_name"`
	Members  []TeamMember `json:"members"`
}

type TeamMember struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type User struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

type OpenPRInfo struct {
	PRID              string   `json:"pr_id"`
	AuthorID          string   `json:"author_id"`
	AssignedReviewers []string `json:"assigned_reviewers"`
}

type PRReviewersUpdate struct {
	PRID      string   `json:"pr_id"`
	Reviewers []string `json:"reviewers"`
}

type MassDeactivationResult struct {
	DeactivatedUsers int      `json:"deactivated_users"`
	UpdatedPRs       int      `json:"updated_prs"`
	FailedPRs        []string `json:"failed_prs,omitempty"`
}

type UserReviewStats struct {
	UserID      string `json:"user_id"`
	Username    string `json:"username"`
	TeamName    string `json:"team_name"`
	ReviewCount int    `json:"review_count"`
}

type PRReviewStats struct {
	PRID          string `json:"pr_id"`
	PRName        string `json:"pr_name"`
	AuthorID      string `json:"author_id"`
	Status        string `json:"status"`
	ReviewerCount int    `json:"reviewer_count"`
}

type OverallStats struct {
	TotalPRs        int     `json:"total_prs"`
	OpenPRs         int     `json:"open_prs"`
	MergedPRs       int     `json:"merged_prs"`
	TotalUsers      int     `json:"total_users"`
	ActiveUsers     int     `json:"active_users"`
	TotalReviews    int     `json:"total_reviews"`
	AvgReviewsPerPR float64 `json:"avg_reviews_per_pr"`
}
