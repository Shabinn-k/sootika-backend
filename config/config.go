package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server struct {
		Port string
	}
	DB struct {
		Host     string
		Port     int
		User     string
		Password string
		Name     string
		SSLMode  string
		TimeZone string
	}
	JWT struct {
		AccessSecret     string
		RefreshSecret    string
		AccessTTLMinutes int
		RefreshTTLHours  int
		MaxSessionHours  int
	}
	SMTP struct {
		Host     string
		Port     int
		Username string
		Password string
		From     string
	}
	OTP struct {
		Length        int
		ExpiryMinutes int
	}
}

func LoadConfig() *Config {
	godotenv.Load()

	cfg := &Config{}

	cfg.Server.Port = getEnv("SERVER_PORT", "8080")

	//database
	cfg.DB.Host = getEnv("DB_HOST", "localhost")
	cfg.DB.Port = getEnvAsInt("DB_PORT", 2007)
	cfg.DB.User = getEnv("DB_USER", "postgres")
	cfg.DB.Password = getEnv("DB_PASSWORD", "")
	cfg.DB.Name = getEnv("DB_NAME", "Sootika")
	cfg.DB.SSLMode = getEnv("DB_SSLMODE", "disable")
	cfg.DB.TimeZone = getEnv("DB_TIMEZONE", "Asia/Kolkata")

	//jwt
	cfg.JWT.AccessSecret = getEnv("ACCESS_SECRET", "")
	cfg.JWT.RefreshSecret = getEnv("REFRESH_SECRET", "")
	cfg.JWT.AccessTTLMinutes = getEnvAsInt("ACCESS_TTL_MINUTE", 15)
	cfg.JWT.RefreshTTLHours = getEnvAsInt("REFRESH_TTL_HOUR", 168)
	cfg.JWT.MaxSessionHours = getEnvAsInt("MAX_SESSION", 720)

	//email
	cfg.SMTP.Host = getEnv("SMTP_HOST", "smtp.gmail.com")
	cfg.SMTP.Port = getEnvAsInt("SMTP_PORT", 587)
	cfg.SMTP.Username = getEnv("SMTP_USERNAME", "")
	cfg.SMTP.Password = getEnv("SMTP_PASSWORD", "")
	cfg.SMTP.From = getEnv("SMTP_FROM", "Sootika <sootika@gmail.com>")

	//otp
	cfg.OTP.Length = getEnvAsInt("OTP_LENGTH", 5)
	cfg.OTP.ExpiryMinutes = getEnvAsInt("OTP_EXPIRY_MINUTES", 5)

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
