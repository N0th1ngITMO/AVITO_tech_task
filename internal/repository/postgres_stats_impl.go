package repository

import (
	"context"
	"database/sql"
	"fmt"
	models "pr_task/internal/model"
)

func (r *PostgresRepository) GetReviewStatsByUser(ctx context.Context) ([]models.UserReviewStats, error) {
	query := `
        SELECT 
            u.user_id,
            u.username, 
            u.team_name,
            (
                SELECT COUNT(*) 
                FROM pull_request pr 
                WHERE u.user_id = ANY(pr.assigned_reviewers)
            ) as review_count
        FROM "user" u
        WHERE u.is_active = true
        ORDER BY review_count DESC
    `

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get user review stats: %v", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Printf("failed to close rows: %v\n", err)
		}
	}(rows)

	var stats []models.UserReviewStats
	for rows.Next() {
		var stat models.UserReviewStats
		if err := rows.Scan(&stat.UserID, &stat.Username, &stat.TeamName, &stat.ReviewCount); err != nil {
			return nil, err
		}
		stats = append(stats, stat)
	}

	return stats, nil
}

func (r *PostgresRepository) GetReviewStatsByPR(ctx context.Context) ([]models.PRReviewStats, error) {
	query := `
		SELECT 
			pull_request_id,
			pull_request_name,
			author_id,
			status,
			COALESCE(array_length(assigned_reviewers, 1), 0) as reviewer_count
		FROM pull_request
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get PR review stats: %v", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Printf("failed to close rows: %v\n", err)
		}
	}(rows)

	var stats []models.PRReviewStats
	for rows.Next() {
		var stat models.PRReviewStats
		if err := rows.Scan(&stat.PRID, &stat.PRName, &stat.AuthorID, &stat.Status, &stat.ReviewerCount); err != nil {
			return nil, err
		}
		stats = append(stats, stat)
	}

	return stats, nil
}

func (r *PostgresRepository) GetOverallStats(ctx context.Context) (*models.OverallStats, error) {
	query := `
		SELECT 
			(SELECT COUNT(*) FROM pull_request) as total_prs,
			(SELECT COUNT(*) FROM pull_request WHERE status = 'OPEN') as open_prs,
			(SELECT COUNT(*) FROM pull_request WHERE status = 'MERGED') as merged_prs,
			(SELECT COUNT(*) FROM "user") as total_users,
			(SELECT COUNT(*) FROM "user" WHERE is_active = true) as active_users,
			(SELECT SUM(COALESCE(array_length(assigned_reviewers, 1), 0)) FROM pull_request) as total_reviews,
			(SELECT AVG(COALESCE(array_length(assigned_reviewers, 1), 0)) FROM pull_request) as avg_reviews
	`

	var stats models.OverallStats
	err := r.db.QueryRowContext(ctx, query).Scan(
		&stats.TotalPRs,
		&stats.OpenPRs,
		&stats.MergedPRs,
		&stats.TotalUsers,
		&stats.ActiveUsers,
		&stats.TotalReviews,
		&stats.AvgReviewsPerPR,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get overall stats: %v", err)
	}

	return &stats, nil
}
