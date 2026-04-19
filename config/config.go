package config
import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
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

	Razorpay struct {
		KeyID     string
		KeySecret string
	}
}

func LoadConfig() *Config {
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load(".env")
		if err != nil {
			log.Println("No .env file found")
		}
	}
	log.Println("ENV DB_HOST:", os.Getenv("DB_HOST"))

	cfg := &Config{}

	cfg.Server.Port = getEnv("SERVER_PORT", "8080")

	// DB
	cfg.DB.Host = mustGetEnv("DB_HOST")
	cfg.DB.Port = getEnvAsInt("DB_PORT", 5432)
	cfg.DB.User = mustGetEnv("DB_USER")
	cfg.DB.Password = mustGetEnv("DB_PASSWORD")
	cfg.DB.Name = mustGetEnv("DB_NAME")
	cfg.DB.SSLMode = getEnv("DB_SSLMODE", "require")
	cfg.DB.TimeZone = getEnv("DB_TIMEZONE", "UTC")

	//jwt
	cfg.JWT.AccessSecret = mustGetEnv("ACCESS_SECRET")
	cfg.JWT.RefreshSecret = mustGetEnv("REFRESH_SECRET")
	cfg.JWT.AccessTTLMinutes = getEnvAsInt("ACCESS_TTL_MINUTE", 15)
	cfg.JWT.RefreshTTLHours = getEnvAsInt("REFRESH_TTL_HOUR", 168)
	cfg.JWT.MaxSessionHours = getEnvAsInt("MAX_SESSION", 720)

	//email
	cfg.SMTP.Host = getEnv("SMTP_HOST", "smtp.gmail.com")
	cfg.SMTP.Port = getEnvAsInt("SMTP_PORT", 587)
	cfg.SMTP.Username = getEnv("SMTP_USERNAME", "")
	cfg.SMTP.Password = getEnv("SMTP_PASSWORD", "")
	cfg.SMTP.From = getEnv("SMTP_FROM", "Sootika <shhabin.07@gmail.com>")

	//otp
	cfg.OTP.Length = getEnvAsInt("OTP_LENGTH", 5)
	cfg.OTP.ExpiryMinutes = getEnvAsInt("OTP_EXPIRY_MINUTES", 5)

	// Razorpay
	cfg.Razorpay.KeyID = mustGetEnv("RAZORPAY_KEY_ID")
	cfg.Razorpay.KeySecret = mustGetEnv("RAZORPAY_KEY_SECRET")
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

func mustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Missing required env: %s", key)
	}
	return value
}
