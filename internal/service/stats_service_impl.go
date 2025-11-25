package services

import (
	"context"
	"pr_task/internal/dto"
)

func (s *ServiceImpl) GetUserReviewStats(ctx context.Context) ([]dto.UserReviewStatsResponse, error) {
	stats, err := s.repo.GetReviewStatsByUser(ctx)
	if err != nil {
		return nil, err
	}

	var response []dto.UserReviewStatsResponse
	for _, stat := range stats {
		response = append(response, dto.UserReviewStatsResponse{
			UserID:      stat.UserID,
			Username:    stat.Username,
			TeamName:    stat.TeamName,
			ReviewCount: stat.ReviewCount,
		})
	}

	return response, nil
}

func (s *ServiceImpl) GetPRReviewStats(ctx context.Context) ([]dto.PRReviewStatsResponse, error) {
	stats, err := s.repo.GetReviewStatsByPR(ctx)
	if err != nil {
		return nil, err
	}

	var response []dto.PRReviewStatsResponse
	for _, stat := range stats {
		response = append(response, dto.PRReviewStatsResponse{
			PRID:          stat.PRID,
			PRName:        stat.PRName,
			AuthorID:      stat.AuthorID,
			Status:        stat.Status,
			ReviewerCount: stat.ReviewerCount,
		})
	}

	return response, nil
}

func (s *ServiceImpl) GetOverallStats(ctx context.Context) (*dto.OverallStatsResponse, error) {
	stats, err := s.repo.GetOverallStats(ctx)
	if err != nil {
		return nil, err
	}

	return &dto.OverallStatsResponse{
		TotalPRs:        stats.TotalPRs,
		OpenPRs:         stats.OpenPRs,
		MergedPRs:       stats.MergedPRs,
		TotalUsers:      stats.TotalUsers,
		ActiveUsers:     stats.ActiveUsers,
		TotalReviews:    stats.TotalReviews,
		AvgReviewsPerPR: stats.AvgReviewsPerPR,
	}, nil
}
