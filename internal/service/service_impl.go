package services

import (
	"context"
	"math/rand"
	"pr_task/internal/dto"
	"time"

	"pr_task/internal/error"
	"pr_task/internal/model"
	"pr_task/internal/repository"
)

type ServiceImpl struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) Service {
	return &ServiceImpl{repo: repo}
}

func (s *ServiceImpl) CreateTeam(ctx context.Context, team models.Team) (*models.Team, error) {
	exists, err := s.repo.TeamExists(ctx, team.TeamName)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.ErrTeamExists
	}

	if err := s.repo.CreateTeam(ctx, team); err != nil {
		return nil, err
	}

	for _, member := range team.Members {
		if err := s.repo.CreateOrUpdateUser(ctx, member, team.TeamName); err != nil {
			return nil, err
		}
	}

	return &team, nil
}

func (s *ServiceImpl) GetTeam(ctx context.Context, teamName string) (*models.Team, error) {
	team, err := s.repo.GetTeam(ctx, teamName)
	if err != nil {
		if err.Error() == "team not found" {
			return nil, errors.ErrNotFound
		}
		return nil, err
	}
	return team, nil
}

func (s *ServiceImpl) SetUserActive(ctx context.Context, userID string, isActive bool) (*models.User, error) {
	user, err := s.repo.GetUser(ctx, userID)
	if err != nil {
		if err.Error() == "user not found" {
			return nil, errors.ErrNotFound
		}
		return nil, err
	}

	if err := s.repo.UpdateUserActive(ctx, userID, isActive); err != nil {
		return nil, err
	}

	user.IsActive = isActive
	return user, nil
}

func (s *ServiceImpl) GetUserReviewPRs(ctx context.Context, userID string) (*dto.UserReviewResponse, error) {
	if _, err := s.repo.GetUser(ctx, userID); err != nil {
		if err.Error() == "user not found" {
			return nil, errors.ErrNotFound
		}
		return nil, err
	}

	prs, err := s.repo.GetPRsByReviewer(ctx, userID)
	if err != nil {
		return nil, err
	}

	var shortPRs []dto.PullRequestShort
	for _, pr := range prs {
		shortPRs = append(shortPRs, dto.PullRequestShort{
			PullRequestID:   pr.PullRequestID,
			PullRequestName: pr.PullRequestName,
			AuthorID:        pr.AuthorID,
			Status:          pr.Status,
		})
	}

	return &dto.UserReviewResponse{
		UserID:       userID,
		PullRequests: shortPRs,
	}, nil
}

func (s *ServiceImpl) CreatePullRequest(ctx context.Context, prID, name, authorID string) (*dto.PullRequest, error) {
	exists, err := s.repo.PRExists(ctx, prID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.ErrPRExists
	}

	author, err := s.repo.GetUser(ctx, authorID)
	if err != nil {
		if err.Error() == "user not found" {
			return nil, errors.ErrNotFound
		}
		return nil, err
	}

	reviewers, err := s.selectReviewers(ctx, author.TeamName, authorID, 2)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	pr := dto.PullRequest{
		PullRequestID:     prID,
		PullRequestName:   name,
		AuthorID:          authorID,
		Status:            "OPEN",
		AssignedReviewers: reviewers,
		CreatedAt:         &now,
	}

	if err := s.repo.CreatePR(ctx, pr); err != nil {
		return nil, err
	}

	return &pr, nil
}

func (s *ServiceImpl) MergePullRequest(ctx context.Context, prID string) (*dto.PullRequest, error) {
	pr, err := s.repo.GetPR(ctx, prID)
	if err != nil {
		if err.Error() == "PR not found" {
			return nil, errors.ErrNotFound
		}
		return nil, err
	}

	if pr.Status == "MERGED" {
		return pr, nil
	}

	now := time.Now()
	if err := s.repo.UpdatePRStatus(ctx, prID, "MERGED", &now); err != nil {
		return nil, err
	}

	pr.Status = "MERGED"
	pr.MergedAt = &now
	return pr, nil
}

func (s *ServiceImpl) ReassignReviewer(ctx context.Context, prID, oldUserID string) (*dto.ReassignResponse, error) {
	pr, err := s.repo.GetPR(ctx, prID)
	if err != nil {
		if err.Error() == "PR not found" {
			return nil, errors.ErrNotFound
		}
		return nil, err
	}

	if pr.Status == "MERGED" {
		return nil, errors.ErrPRMerged
	}

	if !contains(pr.AssignedReviewers, oldUserID) {
		return nil, errors.ErrNotAssigned
	}

	oldReviewer, err := s.repo.GetUser(ctx, oldUserID)
	if err != nil {
		if err.Error() == "user not found" {
			return nil, errors.ErrNotFound
		}
		return nil, err
	}

	excludeIDs := append(pr.AssignedReviewers, pr.AuthorID)
	newReviewer, err := s.repo.GetRandomActiveTeamMember(ctx, oldReviewer.TeamName, excludeIDs)
	if err != nil {
		if err.Error() == "no active team members available" {
			return nil, errors.ErrNoCandidate
		}
		return nil, err
	}

	newReviewers := replaceElement(pr.AssignedReviewers, oldUserID, newReviewer.UserID)
	if err := s.repo.UpdatePRReviewers(ctx, prID, newReviewers); err != nil {
		return nil, err
	}

	pr.AssignedReviewers = newReviewers
	return &dto.ReassignResponse{
		PR:         pr,
		ReplacedBy: newReviewer.UserID,
	}, nil
}

func (s *ServiceImpl) selectReviewers(ctx context.Context, teamName, excludeUserID string, maxReviewers int) ([]string, error) {
	activeMembers, err := s.repo.GetActiveTeamMembers(ctx, teamName, excludeUserID)
	if err != nil {
		return nil, err
	}

	if len(activeMembers) == 0 {
		return []string{}, nil
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	shuffled := make([]models.User, len(activeMembers))
	copy(shuffled, activeMembers)

	rng.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	count := min(len(shuffled), maxReviewers)
	reviewers := make([]string, count)
	for i := 0; i < count; i++ {
		reviewers[i] = shuffled[i].UserID
	}

	return reviewers, nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func replaceElement(slice []string, old, new string) []string {
	result := make([]string, len(slice))
	for i, item := range slice {
		if item == old {
			result[i] = new
		} else {
			result[i] = item
		}
	}
	return result
}
