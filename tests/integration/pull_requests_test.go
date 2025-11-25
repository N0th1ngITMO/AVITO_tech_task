package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPullRequestIntegration(t *testing.T) {
	ctx := context.Background()

	t.Run("CreatePR_Success", func(t *testing.T) {
		clearTestData()
		err := setupTestData(ctx)
		require.NoError(t, err)

		// Создаем PR
		pr, err := testService.CreatePullRequest(ctx, "pr-101", "Feature X", "u1")
		require.NoError(t, err)

		assert.Equal(t, "pr-101", pr.PullRequestID)
		assert.Equal(t, "Feature X", pr.PullRequestName)
		assert.Equal(t, "u1", pr.AuthorID)
		assert.Equal(t, "OPEN", pr.Status)
		assert.NotEmpty(t, pr.AssignedReviewers)
		assert.Len(t, pr.AssignedReviewers, 2)            // Должно быть 2 ревьювера
		assert.NotContains(t, pr.AssignedReviewers, "u1") // Автор не должен быть ревьювером
	})

	t.Run("CreatePR_AuthorNotFound", func(t *testing.T) {
		clearTestData()

		_, err := testService.CreatePullRequest(ctx, "pr-102", "Feature Y", "nonexistent")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("CreatePR_DuplicateID", func(t *testing.T) {
		clearTestData()
		err := setupTestData(ctx)
		require.NoError(t, err)

		// Создаем PR первый раз
		_, err = testService.CreatePullRequest(ctx, "pr-103", "Feature Z", "u1")
		require.NoError(t, err)

		// Пытаемся создать PR с тем же ID
		_, err = testService.CreatePullRequest(ctx, "pr-103", "Feature Z", "u1")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "PR already exists")
	})

	t.Run("MergePR_Success", func(t *testing.T) {
		clearTestData()
		err := setupTestData(ctx)
		require.NoError(t, err)

		// Создаем PR
		pr, err := testService.CreatePullRequest(ctx, "pr-104", "Feature A", "u1")
		require.NoError(t, err)
		assert.Equal(t, "OPEN", pr.Status)

		// Мержим PR
		mergedPR, err := testService.MergePullRequest(ctx, "pr-104")
		require.NoError(t, err)
		assert.Equal(t, "MERGED", mergedPR.Status)
		assert.NotNil(t, mergedPR.MergedAt)

		// Пытаемся мержить еще раз (идемпотентность)
		mergedPR2, err := testService.MergePullRequest(ctx, "pr-104")
		require.NoError(t, err)
		assert.Equal(t, "MERGED", mergedPR2.Status)
	})

	t.Run("MergePR_NotFound", func(t *testing.T) {
		clearTestData()

		_, err := testService.MergePullRequest(ctx, "nonexistent")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("ReassignReviewer_Success", func(t *testing.T) {
		clearTestData()
		err := setupTestData(ctx)
		require.NoError(t, err)

		// Создаем PR
		pr, err := testService.CreatePullRequest(ctx, "pr-105", "Feature B", "u1")
		require.NoError(t, err)
		originalReviewers := pr.AssignedReviewers
		assert.Len(t, originalReviewers, 2)

		// Переназначаем одного ревьювера
		oldReviewer := originalReviewers[0]
		result, err := testService.ReassignReviewer(ctx, "pr-105", oldReviewer)
		require.NoError(t, err)

		assert.NotEqual(t, oldReviewer, result.ReplacedBy)
		assert.NotContains(t, result.PR.AssignedReviewers, oldReviewer)
		assert.Contains(t, result.PR.AssignedReviewers, result.ReplacedBy)
	})

	t.Run("ReassignReviewer_NotAssigned", func(t *testing.T) {
		clearTestData()
		err := setupTestData(ctx)
		require.NoError(t, err)

		// Создаем PR
		_, err = testService.CreatePullRequest(ctx, "pr-106", "Feature C", "u1")
		require.NoError(t, err)

		// Пытаемся переназначить не назначенного ревьювера
		_, err = testService.ReassignReviewer(ctx, "pr-106", "u5") // u5 из другой команды
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not assigned")
	})

	t.Run("ReassignReviewer_MergedPR", func(t *testing.T) {
		clearTestData()
		err := setupTestData(ctx)
		require.NoError(t, err)

		// Создаем и мержим PR
		pr, err := testService.CreatePullRequest(ctx, "pr-107", "Feature D", "u1")
		require.NoError(t, err)

		_, err = testService.MergePullRequest(ctx, "pr-107")
		require.NoError(t, err)

		// Пытаемся переназначить ревьювера в замерженном PR
		_, err = testService.ReassignReviewer(ctx, "pr-107", pr.AssignedReviewers[0])
		require.Error(t, err)
		assert.Contains(t, err.Error(), "PR is merged")
	})
}
