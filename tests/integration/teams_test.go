package integration

import (
	"context"
	models "pr_task/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTeamIntegration(t *testing.T) {
	ctx := context.Background()

	t.Run("CreateTeam_Success", func(t *testing.T) {
		// Очищаем данные перед тестом
		clearTestData()

		team := models.Team{
			TeamName: "devops",
			Members: []models.TeamMember{
				{UserID: "u10", Username: "Grace", IsActive: true},
				{UserID: "u11", Username: "Henry", IsActive: true},
			},
		}

		// Создаем команду через сервис
		createdTeam, err := testService.CreateTeam(ctx, team)
		require.NoError(t, err)
		assert.Equal(t, team.TeamName, createdTeam.TeamName)
		assert.Len(t, createdTeam.Members, 2)

		// Проверяем, что команда сохранилась в БД
		dbTeam, err := testService.GetTeam(ctx, "devops")
		require.NoError(t, err)
		assert.Equal(t, team.TeamName, dbTeam.TeamName)
	})

	t.Run("CreateTeam_AlreadyExists", func(t *testing.T) {
		// Очищаем данные перед тестом
		clearTestData()

		team := models.Team{
			TeamName: "mobile",
			Members: []models.TeamMember{
				{UserID: "u12", Username: "Ivan", IsActive: true},
			},
		}

		// Создаем команду первый раз
		_, err := testService.CreateTeam(ctx, team)
		require.NoError(t, err)

		// Пытаемся создать команду с тем же именем
		_, err = testService.CreateTeam(ctx, team)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "team already exists")
	})

	t.Run("GetTeam_Success", func(t *testing.T) {
		// Очищаем данные перед тестом
		clearTestData()

		// Создаем тестовые данные
		err := setupTestData(ctx)
		require.NoError(t, err)

		// Получаем команду
		team, err := testService.GetTeam(ctx, "backend")
		require.NoError(t, err)
		assert.Equal(t, "backend", team.TeamName)
		assert.Len(t, team.Members, 4) // 3 активных + 1 неактивный
	})

	t.Run("GetTeam_NotFound", func(t *testing.T) {
		// Очищаем данные перед тестом
		clearTestData()

		_, err := testService.GetTeam(ctx, "nonexistent")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}
