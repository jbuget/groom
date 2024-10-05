package config

import (
	"log"
	"os"
)

type Config struct {
	Host                                 string
	Port                                 string
	DatabaseURL                          string
	DatabaseMigrationPath                string
	Username                             string
	Password                             string
	APIKey                               string
	GoogleWorkspaceDomain                string
	GoogleAPIKey                         string
	GoogleClientID                       string
	GoogleClientSecret                   string
	GoogleRedirectURL                    string
	GoogleServiceAccountCredentialsFile  string
	GoogleServiceAccountImpersonatedUser string // email
}

func LoadConfig() Config {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	googleWorkspaceDomain := os.Getenv("GOOGLE_WORKSPACE_DOMAIN")
	if googleWorkspaceDomain == "" {
		log.Fatal("GOOGLE_WORKSPACE_DOMAIN is not set")
	}

	googleAPIKey := os.Getenv("GOOGLE_API_KEY")
	if googleAPIKey == "" {
		log.Fatal("GOOGLE_API_KEY is not set")
	}

	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	if googleClientID == "" {
		log.Fatal("GOOGLE_CLIENT_ID is not set")
	}

	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	if googleClientSecret == "" {
		log.Fatal("GOOGLE_CLIENT_SECRET is not set")
	}

	googleRedirectURL := os.Getenv("GOOGLE_REDIRECT_URL")
	if googleRedirectURL == "" {
		log.Fatal("GOOGLE_REDIRECT_URL is not set")
	}

	googleServiceAccountImpersonatedUser := os.Getenv("GOOGLE_SERVICE_ACCOUNT_IMPERSONATED_USER")
	if googleServiceAccountImpersonatedUser == "" {
		log.Fatal("GOOGLE_SERVICE_ACCOUNT_IMPERSONATED_USER is not set")
	}

	return Config{
		Host:                                 getEnv("HOST", "0.0.0.0"),
		Port:                                 getEnv("PORT", "3000"),
		DatabaseURL:                          os.Getenv("DATABASE_URL"),
		DatabaseMigrationPath:                getEnv("DATABASE_MIGRATION_PATH", "./migrations"),
		APIKey:                               os.Getenv("GROOM_API_KEY"),
		GoogleWorkspaceDomain:                os.Getenv("GOOGLE_WORKSPACE_DOMAIN"),
		GoogleAPIKey:                         os.Getenv("GOOGLE_API_KEY"),
		GoogleClientID:                       os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret:                   os.Getenv("GOOGLE_CLIENT_SECRET"),
		GoogleRedirectURL:                    os.Getenv("GOOGLE_REDIRECT_URL"),
		GoogleServiceAccountCredentialsFile:  getEnv("GOOGLE_SERVICE_ACCOUNT_CREDENTIALS_FILE", "./service_account.json"),
		GoogleServiceAccountImpersonatedUser: os.Getenv("GOOGLE_SERVICE_ACCOUNT_IMPERSONATED_USER"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
