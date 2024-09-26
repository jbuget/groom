package config

import (
    "os"
)

type Config struct {
    Host        string
    Port        string
    DatabaseURL string
    Username    string
    Password    string
    APIKey      string
}

func LoadConfig() Config {
    return Config{
        Host:        getEnv("HOST", "0.0.0.0"),
        Port:        getEnv("PORT", "3000"),
        DatabaseURL: os.Getenv("DATABASE_URL"),
        Username:    os.Getenv("BASIC_AUTH_LOGIN"),
        Password:    os.Getenv("BASIC_AUTH_PASSWORD"),
        APIKey:      os.Getenv("USHR_API_KEY"),
    }
}

func getEnv(key, defaultValue string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return defaultValue
}