package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"golang/config"
	"time"
)

type Manager struct {
	accessSercret string
	refreshSecret string
	accessTTL     time.Duration
	refreshTTL    time.Duration
	maxSession    time.Duration
}

func NewJWTManager(cfg *config.Config) *Manager {
	return &Manager{
		accessSercret: cfg.JWT.AccessSecret,
		refreshSecret: cfg.JWT.RefreshSecret,
		accessTTL:     time.Duration(cfg.JWT.AccessTTLMinutes) * time.Minute,
		refreshTTL:    time.Duration(cfg.JWT.RefreshTTLHours) * time.Hour,
		maxSession:    time.Duration(cfg.JWT.MaxSessionHours) * time.Hour,
	}
}

func (j *Manager) GenerateAccessToken(userId, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userId,
		"role":    role,
		"exp":     time.Now().Add(j.accessTTL).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.accessSercret))
}

func (j *Manager) GenerateRefreshToken(userId, role, sessionId string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":    userId,
		"role":       role,
		"session_id": sessionId,
		"exp":        time.Now().Add(j.refreshTTL).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.refreshSecret))
}

func (j *Manager) ValidateAccessToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		return []byte(j.accessSercret), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid access token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}
	return claims, nil
}

func (j *Manager) ValidateRefresh(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		return []byte(j.refreshSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}
	return claims, nil
}