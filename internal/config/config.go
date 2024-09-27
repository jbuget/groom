package config

import (
    "log"
    "os"
)

type Config struct {
    Host        string
    Port        string
    DatabaseURL string
    Username    string
    Password    string
    APIKey      string
	GoogleAPIKey	   string
    GoogleClientID     string
    GoogleClientSecret string
    GoogleRedirectURL  string
}

func LoadConfig() Config {
    dbURL := os.Getenv("DATABASE_URL")
    if dbURL == "" {
        log.Fatal("DATABASE_URL is not set")
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

	return Config{
        Host:        getEnv("HOST", "0.0.0.0"),
        Port:        getEnv("PORT", "3000"),
        DatabaseURL: os.Getenv("DATABASE_URL"),
        Username:    os.Getenv("BASIC_AUTH_LOGIN"),
        Password:    os.Getenv("BASIC_AUTH_PASSWORD"),
        APIKey:      os.Getenv("GROOM_API_KEY"),
        GoogleAPIKey:     os.Getenv("GOOGLE_API_KEY"),
        GoogleClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
        GoogleClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
        GoogleRedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
    }
}

func getEnv(key, defaultValue string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return defaultValue
}