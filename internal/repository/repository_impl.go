package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"pr_task/internal/dto"
	"strings"
	"time"

	"pr_task/internal/model"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) Repository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) CreateTeam(ctx context.Context, team models.Team) error {
	query := `INSERT INTO team (team_name) VALUES ($1)`
	_, err := r.db.ExecContext(ctx, query, team.TeamName)
	return err
}

func (r *PostgresRepository) GetTeam(ctx context.Context, teamName string) (*models.Team, error) {
	var team models.Team
	query := `SELECT team_name FROM team WHERE team_name = $1`
	err := r.db.QueryRowContext(ctx, query, teamName).Scan(&team.TeamName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("team not found")
		}
		return nil, err
	}

	membersQuery := `SELECT user_id, username, is_active FROM "user" WHERE team_name = $1`
	rows, err := r.db.QueryContext(ctx, membersQuery, teamName)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(rows)

	var members []models.TeamMember
	for rows.Next() {
		var member models.TeamMember
		if err := rows.Scan(&member.UserID, &member.Username, &member.IsActive); err != nil {
			return nil, err
		}
		members = append(members, member)
	}

	team.Members = members
	return &team, nil
}

func (r *PostgresRepository) TeamExists(ctx context.Context, teamName string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM team WHERE team_name = $1)`
	err := r.db.QueryRowContext(ctx, query, teamName).Scan(&exists)
	return exists, err
}

func (r *PostgresRepository) CreateOrUpdateUser(ctx context.Context, user models.TeamMember, teamName string) error {
	query := `
		INSERT INTO "user" (user_id, username, team_name, is_active) 
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id) 
		DO UPDATE SET username = $2, team_name = $3, is_active = $4
	`
	_, err := r.db.ExecContext(ctx, query, user.UserID, user.Username, teamName, user.IsActive)
	return err
}

func (r *PostgresRepository) GetUser(ctx context.Context, userID string) (*models.User, error) {
	var user models.User
	query := `SELECT user_id, username, team_name, is_active FROM "user" WHERE user_id = $1`
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&user.UserID, &user.Username, &user.TeamName, &user.IsActive)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *PostgresRepository) UpdateUserActive(ctx context.Context, userID string, isActive bool) error {
	query := `UPDATE "user" SET is_active = $1 WHERE user_id = $2`
	result, err := r.db.ExecContext(ctx, query, isActive, userID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

func (r *PostgresRepository) GetActiveTeamMembers(ctx context.Context, teamName string, excludeUserID string) ([]models.User, error) {
	query := `SELECT user_id, username, team_name, is_active FROM "user" WHERE team_name = $1 AND is_active = true AND user_id != $2`
	rows, err := r.db.QueryContext(ctx, query, teamName, excludeUserID)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(rows)

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.UserID, &user.Username, &user.TeamName, &user.IsActive); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *PostgresRepository) GetRandomActiveTeamMember(ctx context.Context, teamName string, excludeUserIDs []string) (*models.User, error) {
	var placeholders []string
	var args []interface{}
	args = append(args, teamName)

	for i, userID := range excludeUserIDs {
		placeholders = append(placeholders, fmt.Sprintf("$%d", i+2))
		args = append(args, userID)
	}

	excludeClause := ""
	if len(excludeUserIDs) > 0 {
		excludeClause = fmt.Sprintf("AND user_id NOT IN (%s)", strings.Join(placeholders, ", "))
	}

	query := `
    SELECT user_id, username, team_name, is_active 
    FROM "user" 
    WHERE team_name = $1 AND is_active = true ` + excludeClause + `
    ORDER BY RANDOM()
    LIMIT 1
`

	var user models.User
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&user.UserID, &user.Username, &user.TeamName, &user.IsActive)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no active team members available")
		}
		return nil, err
	}
	return &user, nil
}

func (r *PostgresRepository) CreatePR(ctx context.Context, pr dto.PullRequest) error {
	query := `
		INSERT INTO pull_request (pull_request_id, pull_request_name, author_id, status, assigned_reviewers, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctx, query,
		pr.PullRequestID,
		pr.PullRequestName,
		pr.AuthorID,
		pr.Status,
		pq.Array(pr.AssignedReviewers),
		pr.CreatedAt,
	)
	return err
}

func (r *PostgresRepository) GetPR(ctx context.Context, prID string) (*dto.PullRequest, error) {
	var pr dto.PullRequest
	query := `
		SELECT pull_request_id, pull_request_name, author_id, status, assigned_reviewers, created_at, merged_at
		FROM pull_request 
		WHERE pull_request_id = $1
	`

	var reviewers []string
	err := r.db.QueryRowContext(ctx, query, prID).Scan(
		&pr.PullRequestID,
		&pr.PullRequestName,
		&pr.AuthorID,
		&pr.Status,
		pq.Array(&reviewers),
		&pr.CreatedAt,
		&pr.MergedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("PR not found")
		}
		return nil, err
	}

	pr.AssignedReviewers = reviewers
	return &pr, nil
}

func (r *PostgresRepository) PRExists(ctx context.Context, prID string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM pull_request WHERE pull_request_id = $1)`
	err := r.db.QueryRowContext(ctx, query, prID).Scan(&exists)
	return exists, err
}

func (r *PostgresRepository) UpdatePRStatus(ctx context.Context, prID string, status string, mergedAt *time.Time) error {
	query := `UPDATE pull_request SET status = $1, merged_at = $2 WHERE pull_request_id = $3`
	result, err := r.db.ExecContext(ctx, query, status, mergedAt, prID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("PR not found")
	}
	return nil
}

func (r *PostgresRepository) UpdatePRReviewers(ctx context.Context, prID string, reviewers []string) error {
	query := `UPDATE pull_request SET assigned_reviewers = $1 WHERE pull_request_id = $2`
	result, err := r.db.ExecContext(ctx, query, pq.Array(reviewers), prID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("PR not found")
	}
	return nil
}

func (r *PostgresRepository) GetPRsByReviewer(ctx context.Context, userID string) ([]dto.PullRequest, error) {
	query := `
		SELECT pull_request_id, pull_request_name, author_id, status, assigned_reviewers, created_at, merged_at
		FROM pull_request 
		WHERE $1 = ANY(assigned_reviewers)
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(rows)

	var prs []dto.PullRequest
	for rows.Next() {
		var pr dto.PullRequest
		var reviewers []string
		err := rows.Scan(
			&pr.PullRequestID,
			&pr.PullRequestName,
			&pr.AuthorID,
			&pr.Status,
			pq.Array(&reviewers),
			&pr.CreatedAt,
			&pr.MergedAt,
		)
		if err != nil {
			return nil, err
		}
		pr.AssignedReviewers = reviewers
		prs = append(prs, pr)
	}
	return prs, nil
}
