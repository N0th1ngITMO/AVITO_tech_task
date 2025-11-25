package repository

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
