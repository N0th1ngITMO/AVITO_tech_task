package repository

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
