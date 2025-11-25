package integration

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	handlers "pr_task/internal/handler"
	"pr_task/internal/repository"
	services "pr_task/internal/service"
	"testing"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var (
	testDB      *sql.DB
	testRepo    repository.Repository
	testService services.Service
	testHandler *handlers.Handler
)

// TestMain настраивает тестовое окружение
func TestMain(m *testing.M) {
	if err := setup(); err != nil {
		log.Fatalf("Failed to setup tests: %v", err)
	}

	code := m.Run()

	cleanup()
	os.Exit(code)
}

// setup инициализирует тестовую базу данных и сервисы
func setup() error {
	// Загружаем .env.test файл
	if err := godotenv.Load("../../.env.test"); err != nil {
		log.Println("No .env.test file found, using default values")
	}

	// Параметры подключения к тестовой БД
	dbHost := getEnv("TEST_DB_HOST", "localhost")
	dbPort := getEnv("TEST_DB_PORT", "5433")
	dbUser := getEnv("TEST_DB_USER", "test_user")
	dbPassword := getEnv("TEST_DB_PASSWORD", "test_pass")
	dbName := getEnv("TEST_DB_NAME", "pr_review_test")
	dbSSLMode := getEnv("TEST_DB_SSLMODE", "disable")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode)

	// Подключаемся к базе данных
	var err error
	testDB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to connect to test database: %v", err)
	}

	// Проверяем подключение
	if err := testDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping test database: %v", err)
	}

	// Создаем таблицы из init.sql
	if err := createTestTables(); err != nil {
		return fmt.Errorf("failed to create test tables: %v", err)
	}

	// Инициализируем зависимости
	testRepo = repository.NewPostgresRepository(testDB)
	testService = services.NewService(testRepo)
	testHandler = handlers.NewHandler(testService)

	log.Println("Test database setup completed successfully")
	return nil
}

// createTestTables создает таблицы для тестов
func createTestTables() error {
	queries := []string{
		// Создаем ENUM тип если не существует
		`DO $$ BEGIN
			CREATE TYPE pr_status AS ENUM ('OPEN', 'MERGED');
		EXCEPTION
			WHEN duplicate_object THEN null;
		END $$;`,

		// Таблица команд
		`CREATE TABLE IF NOT EXISTS team (
			team_name TEXT PRIMARY KEY
		)`,

		// Таблица пользователей
		`CREATE TABLE IF NOT EXISTS "user" (
			user_id   TEXT PRIMARY KEY,
			username  TEXT NOT NULL,
			team_name TEXT NOT NULL REFERENCES team(team_name) ON DELETE CASCADE,
			is_active BOOLEAN NOT NULL DEFAULT true
		)`,

		// Таблица Pull Requests
		`CREATE TABLE IF NOT EXISTS pull_request (
			pull_request_id   TEXT PRIMARY KEY,
			pull_request_name TEXT NOT NULL,
			author_id         TEXT NOT NULL REFERENCES "user"(user_id),
			status            pr_status NOT NULL DEFAULT 'OPEN',
			assigned_reviewers TEXT[] NOT NULL DEFAULT '{}',
			created_at        TIMESTAMPTZ DEFAULT NOW(),
			merged_at         TIMESTAMPTZ
		)`,

		// Создаем индексы
		`CREATE INDEX IF NOT EXISTS idx_user_team_name ON "user"(team_name)`,
		`CREATE INDEX IF NOT EXISTS idx_pr_author_id ON pull_request(author_id)`,
		`CREATE INDEX IF NOT EXISTS idx_pr_status ON pull_request(status)`,
		`CREATE INDEX IF NOT EXISTS idx_pr_reviewers ON pull_request USING GIN (assigned_reviewers)`,
	}

	for _, query := range queries {
		if _, err := testDB.Exec(query); err != nil {
			return fmt.Errorf("failed to execute query: %s, error: %v", query, err)
		}
	}

	return nil
}

// cleanup очищает тестовую базу данных
func cleanup() {
	if testDB != nil {
		clearTestData()
		testDB.Close()
	}
}

// clearTestData очищает тестовые данные
func clearTestData() {
	queries := []string{
		"DELETE FROM pull_request",
		"DELETE FROM \"user\"",
		"DELETE FROM team",
	}

	for _, query := range queries {
		testDB.Exec(query)
	}
}

// getEnv возвращает переменную окружения или значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// setupTestData создает тестовые данные
func setupTestData(ctx context.Context) error {
	// Создаем команды
	teams := []struct {
		name    string
		members []struct {
			userID   string
			username string
			active   bool
		}
	}{
		{
			name: "backend",
			members: []struct {
				userID   string
				username string
				active   bool
			}{
				{"u1", "Alice", true},
				{"u2", "Bob", true},
				{"u3", "Charlie", true},
				{"u4", "David", false}, // неактивный
			},
		},
		{
			name: "frontend",
			members: []struct {
				userID   string
				username string
				active   bool
			}{
				{"u5", "Eve", true},
				{"u6", "Frank", true},
			},
		},
		{
			name: "devops",
			members: []struct {
				userID   string
				username string
				active   bool
			}{
				{"u7", "Grace", true},
				{"u8", "Henry", true},
			},
		},
	}

	for _, team := range teams {
		// Создаем команду
		if _, err := testDB.ExecContext(ctx, "INSERT INTO team (team_name) VALUES ($1) ON CONFLICT DO NOTHING", team.name); err != nil {
			return err
		}

		// Создаем пользователей
		for _, member := range team.members {
			_, err := testDB.ExecContext(ctx,
				"INSERT INTO \"user\" (user_id, username, team_name, is_active) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING",
				member.userID, member.username, team.name, member.active)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
