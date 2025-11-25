package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserIntegration(t *testing.T) {
	ctx := context.Background()

	t.Run("SetUserActive_Success", func(t *testing.T) {
		// Очищаем и настраиваем данные
		clearTestData()
		err := setupTestData(ctx)
		require.NoError(t, err)

		// Деактивируем пользователя
		user, err := testService.SetUserActive(ctx, "u1", false)
		require.NoError(t, err)
		assert.False(t, user.IsActive)
		assert.Equal(t, "u1", user.UserID)

		// Активируем обратно
		user, err = testService.SetUserActive(ctx, "u1", true)
		require.NoError(t, err)
		assert.True(t, user.IsActive)
	})

	t.Run("SetUserActive_UserNotFound", func(t *testing.T) {
		clearTestData()

		_, err := testService.SetUserActive(ctx, "nonexistent", true)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("GetUserReviewPRs_Success", func(t *testing.T) {
		clearTestData()
		err := setupTestData(ctx)
		require.NoError(t, err)

		// Создаем PR с ревьюверами
		pr, err := testService.CreatePullRequest(ctx, "pr-100", "Test PR", "u1")
		require.NoError(t, err)
		assert.Contains(t, pr.AssignedReviewers, "u2") // u2 должен быть назначен ревьювером

		// Получаем PR для ревью пользователя u2
		result, err := testService.GetUserReviewPRs(ctx, "u2")
		require.NoError(t, err)
		assert.Equal(t, "u2", result.UserID)
		assert.Len(t, result.PullRequests, 1)
		assert.Equal(t, "pr-100", result.PullRequests[0].PullRequestID)
	})

	t.Run("GetUserReviewPRs_NoReviews", func(t *testing.T) {
		clearTestData()
		err := setupTestData(ctx)
		require.NoError(t, err)

		// Пользователь u5 не должен иметь PR для ревью
		result, err := testService.GetUserReviewPRs(ctx, "u5")
		require.NoError(t, err)
		assert.Equal(t, "u5", result.UserID)
		assert.Empty(t, result.PullRequests)
	})

	t.Run("GetUserReviewPRs_UserNotFound", func(t *testing.T) {
		clearTestData()

		_, err := testService.GetUserReviewPRs(ctx, "nonexistent")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}
