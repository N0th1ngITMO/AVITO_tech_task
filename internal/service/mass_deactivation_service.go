package services

import (
	"context"
	"fmt"
	"pr_task/internal/repository"
	"time"
)

type MassDeactivationResponse struct {
	DeactivatedUsers int      `json:"deactivated_users" example:"5"`
	UpdatedPRs       int      `json:"updated_prs" example:"3"`
	FailedPRs        []string `json:"failed_prs,omitempty" example:"pr-1001"`
	ProcessingTime   int64    `json:"processing_time_ms" example:"85"`
}

func (s *ServiceImpl) MassDeactivateTeamUsers(ctx context.Context, teamName string, excludeUserIDs []string) (*MassDeactivationResponse, error) {
	startTime := time.Now()

	_, err := s.repo.GetTeam(ctx, teamName)
	if err != nil {
		return nil, fmt.Errorf("team not found: %v", err)
	}

	deactivatedCount, err := s.repo.MassDeactivateUsers(ctx, teamName, excludeUserIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to deactivate users: %v", err)
	}

	if deactivatedCount == 0 {
		return &MassDeactivationResponse{
			DeactivatedUsers: 0,
			UpdatedPRs:       0,
			FailedPRs:        []string{},
			ProcessingTime:   time.Since(startTime).Milliseconds(),
		}, nil
	}

	openPRs, err := s.repo.GetOpenPRsWithReviewers(ctx, teamName)
	if err != nil {
		return &MassDeactivationResponse{
			DeactivatedUsers: deactivatedCount,
			UpdatedPRs:       0,
			FailedPRs:        []string{"all"},
			ProcessingTime:   time.Since(startTime).Milliseconds(),
		}, nil
	}

	updateResult := s.updateReviewersForOpenPRs(ctx, openPRs, excludeUserIDs)

	processingTime := time.Since(startTime).Milliseconds()

	return &MassDeactivationResponse{
		DeactivatedUsers: deactivatedCount,
		UpdatedPRs:       updateResult.UpdatedPRs,
		FailedPRs:        updateResult.FailedPRs,
		ProcessingTime:   processingTime,
	}, nil
}

func (s *ServiceImpl) updateReviewersForOpenPRs(ctx context.Context, openPRs []repository.OpenPRInfo, excludeUserIDs []string) *repository.MassDeactivationResult {
	if len(openPRs) == 0 {
		return &repository.MassDeactivationResult{UpdatedPRs: 0}
	}

	var updates []repository.PRReviewersUpdate
	var failedPRs []string

	for _, pr := range openPRs {
		newReviewers := s.getUpdatedReviewers(pr.AssignedReviewers, excludeUserIDs, pr.AuthorID)

		updates = append(updates, repository.PRReviewersUpdate{
			PRID:      pr.PRID,
			Reviewers: newReviewers,
		})
	}

	if len(updates) > 0 {
		if err := s.repo.UpdatePRReviewersBatch(ctx, updates); err != nil {
			for _, update := range updates {
				failedPRs = append(failedPRs, update.PRID)
			}
			return &repository.MassDeactivationResult{
				UpdatedPRs: 0,
				FailedPRs:  failedPRs,
			}
		}
	}

	return &repository.MassDeactivationResult{
		UpdatedPRs: len(updates),
		FailedPRs:  failedPRs,
	}
}

func (s *ServiceImpl) getUpdatedReviewers(currentReviewers []string, availableUsers []string, authorID string) []string {
	newReviewers := make([]string, 0, 2)

	for _, reviewer := range currentReviewers {
		if reviewer == authorID {
			continue
		}

		if contains(availableUsers, reviewer) {
			newReviewers = append(newReviewers, reviewer)
			continue
		}

		replacement := s.findAvailableReviewer(availableUsers, newReviewers, authorID)

		fmt.Println(replacement)
		if replacement != "" {
			newReviewers = append(newReviewers, replacement)
		}
	}

	if len(newReviewers) > 2 {
		newReviewers = newReviewers[:2]
	}

	return newReviewers
}

func (s *ServiceImpl) findAvailableReviewer(availableUsers []string, currentReviewers []string, authorID string) string {
	excludeUsers := append(currentReviewers, authorID)

	for _, user := range availableUsers {
		if !contains(excludeUsers, user) {
			return user
		}
	}
	return ""
}
