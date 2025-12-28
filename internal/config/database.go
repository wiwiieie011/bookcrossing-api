package config

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(logger *slog.Logger) *gorm.DB {
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	dbMode := os.Getenv("DB_SSLMODE")

	if dbMode == "" {
		dbMode = "disable"
	}

	// Add short connect_timeout to fail-fast per attempt
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=%v connect_timeout=%d",
		dbHost, dbUser, dbPass, dbName, dbPort, dbMode, 5)

	// Retry connection with backoff to handle DNS/startup races
	var (
		db  *gorm.DB
		err error
	)
	maxAttempts := 12 // ~1m total
	backoff := 2 * time.Second
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		logger.Warn("database connect attempt failed", "attempt", attempt, "error", err)
		time.Sleep(backoff)
		// Exponential backoff but cap at 10s
		if backoff < 10*time.Second {
			backoff *= 2
			if backoff > 10*time.Second {
				backoff = 10 * time.Second
			}
		}
	}
	if err != nil {
		logger.Error("failed to connect to database after retries", "error", err)
		os.Exit(1)
	}

	// Получаем *sql.DB для настройки пула соединений
	sqlDB, err := db.DB()
	if err != nil {
		logger.Error("failed to get sql.DB from gorm.DB", "error", err)
		os.Exit(1)
	}

	// Настройка пула соединений
	sqlDB.SetMaxOpenConns(50)                 // максимум открытых соединений
	sqlDB.SetMaxIdleConns(25)                 // сколько соединений может простаивать
	sqlDB.SetConnMaxLifetime(time.Minute * 5) // максимальное время жизни соединения

	logger.Info("connected to database with connection pool configured")

	return db
}
