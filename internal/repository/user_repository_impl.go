package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	models "pr_task/internal/model"
)

func (r *PostgresRepository) MassDeactivateUsers(ctx context.Context, teamName string, excludeUserIDs []string) (int, error) {
	var query string
	var args []interface{}

	if len(excludeUserIDs) == 0 {
		query = `UPDATE "user" SET is_active = false WHERE team_name = $1 AND is_active = true`
		args = []interface{}{teamName}
	} else {
		query = `UPDATE "user" SET is_active = false WHERE team_name = $1 AND is_active = true AND NOT (user_id = ANY($2))`
		args = []interface{}{teamName, pq.Array(excludeUserIDs)}
	}

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to mass deactivate users: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %v", err)
	}

	return int(rowsAffected), nil
}

func (r *PostgresRepository) GetOpenPRsWithReviewers(ctx context.Context, teamName string) ([]models.OpenPRInfo, error) {
	query := `
		SELECT 
			pr.pull_request_id,
			pr.assigned_reviewers,
			u.user_id as author_id
		FROM pull_request pr
		JOIN "user" u ON pr.author_id = u.user_id
		WHERE pr.status = 'OPEN' 
		AND u.team_name = $1
	`

	rows, err := r.db.QueryContext(ctx, query, teamName)
	if err != nil {
		return nil, fmt.Errorf("failed to get open PRs: %v", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Printf("failed to close rows: %v", err)
		}
	}(rows)

	var prs []models.OpenPRInfo
	for rows.Next() {
		var pr models.OpenPRInfo
		var reviewers []string

		if err := rows.Scan(&pr.PRID, pq.Array(&reviewers), &pr.AuthorID); err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}

		pr.AssignedReviewers = reviewers
		prs = append(prs, pr)
	}

	return prs, nil
}

func (r *PostgresRepository) UpdatePRReviewersBatch(ctx context.Context, updates []models.PRReviewersUpdate) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer func(tx *sql.Tx) {
		err := tx.Rollback()
		if err != nil {
			fmt.Printf("failed to rollback transaction: %v", err)
		}
	}(tx)

	stmt, err := tx.PrepareContext(ctx, `
		UPDATE pull_request 
		SET assigned_reviewers = $1 
		WHERE pull_request_id = $2
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			fmt.Printf("failed to close statement: %v", err)
		}
	}(stmt)

	for _, update := range updates {
		if _, err := stmt.ExecContext(ctx, pq.Array(update.Reviewers), update.PRID); err != nil {
			return fmt.Errorf("failed to update PR %s: %v", update.PRID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}
