package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	_ "pr_task/docs"
	"pr_task/internal/config"
	handlers "pr_task/internal/handler"
	"pr_task/internal/repository"
	"pr_task/internal/routes"
	services "pr_task/internal/service"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func LoadConfig() (*config.DB, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found, using environment variables")
	}

	configDB := &config.DB{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "pr_user"),
		DBPassword: getEnv("DB_PASSWORD", "secure_password_123"),
		DBName:     getEnv("DB_NAME", "pr_review_db"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
	}

	maxOpenConns, err := strconv.Atoi(getEnv("DB_MAX_OPEN_CONNS", "25"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_MAX_OPEN_CONNS: %v", err)
	}
	configDB.DBMaxOpenConns = maxOpenConns

	maxIdleConns, err := strconv.Atoi(getEnv("DB_MAX_IDLE_CONNS", "25"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_MAX_IDLE_CONNS: %v", err)
	}
	configDB.DBMaxIdleConns = maxIdleConns

	connMaxLifetime, err := time.ParseDuration(getEnv("DB_CONN_MAX_LIFETIME", "5m"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_CONN_MAX_LIFETIME: %v", err)
	}
	configDB.DBConnMaxLifetime = connMaxLifetime

	return configDB, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func createDBConnection(config *config.DB) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName, config.DBSSLMode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %v", err)
	}

	db.SetMaxOpenConns(config.DBMaxOpenConns)
	db.SetMaxIdleConns(config.DBMaxIdleConns)
	db.SetConnMaxLifetime(config.DBConnMaxLifetime)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	log.Printf("Successfully connected to database: %s@%s:%s/%s",
		config.DBUser, config.DBHost, config.DBPort, config.DBName)

	return db, nil
}

// @title PR Reviewer Assignment Service
// @version 1.0.0
// @description Сервис автоматического назначения ревьюверов для Pull Request'ов

// @host localhost:8080
// @BasePath /

func main() {
	configDB, err := LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	db, err := createDBConnection(configDB)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			fmt.Printf("Failed to close database connection: %v", err)
		}
	}(db)

	e := echo.New()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().Method == "POST" || c.Request().Method == "PUT" {
				log.Printf("Request to %s", c.Path())
			}
			return next(c)
		}
	})

	repo := repository.NewPostgresRepository(db)
	service := services.NewService(repo)
	handler := handlers.NewHandler(service)

	routes.RegisterRoutes(e, handler)

	serverAddress := ":" + configDB.ServerPort
	log.Printf("Server starting on %s", serverAddress)
	e.Logger.Fatal(e.Start(serverAddress))
}
