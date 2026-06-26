package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	AppUrl             string
	AppPort            string
	DBHost             string
	DBPort             string
	DBName             string
	DBUser             string
	DBPassword         string
	SMTPHost           string
	SMTPPort           string
	SMTPUsername       string
	SMTPPassword       string
	FrontendURL        string
	FrontendSuccessURL string
	FrontendErrorURL   string
	WebinarDateTime    string
	WebinarLink        string
}

var AppConfig Config

func Load() {
	_ = godotenv.Load()

	AppConfig = Config{
		AppUrl:             getEnv("APP_URL", "localhost"),
		AppPort:            getEnv("APP_PORT", "8080"),
		DBHost:             getEnv("DB_HOST", "localhost"),
		DBPort:             getEnv("DB_PORT", "3306"),
		DBName:             getEnv("DB_NAME", "app_db"),
		DBUser:             getEnv("DB_USER", "root"),
		DBPassword:         getEnv("DB_PASSWORD", "password"),
		SMTPHost:           getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:           getEnv("SMTP_PORT", "587"),
		SMTPUsername:       getEnv("SMTP_USERNAME", ""),
		SMTPPassword:       getEnv("SMTP_PASSWORD", ""),
		FrontendURL:        getEnv("FRONTEND_URL", "http://localhost:3000"),
		FrontendSuccessURL: getEnv("FRONTEND_SUCCESS_URL", "http://localhost:3000/success"),
		FrontendErrorURL:   getEnv("FRONTEND_ERROR_URL", "http://localhost:3000/error"),
		WebinarDateTime:    getEnv("WEBINAR_DATETIME", "TBD"),
		WebinarLink:        getEnv("WEBINAR_LINK", "https://example.com/webinar"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func ConnectDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		AppConfig.DBUser,
		AppConfig.DBPassword,
		AppConfig.DBHost,
		AppConfig.DBPort,
		AppConfig.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}
