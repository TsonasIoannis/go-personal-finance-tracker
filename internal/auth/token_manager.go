package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/models"
)

const defaultJWTSecret = "development-secret-change-me"

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

type Claims struct {
	UserID    uint   `json:"user_id"`
	Email     string `json:"email"`
	ExpiresAt int64  `json:"exp"`
}

type TokenManager interface {
	GenerateToken(user *models.User) (string, error)
	ParseToken(token string) (*Claims, error)
}

type JWTManager struct {
	secret []byte
	ttl    time.Duration
}

func NewJWTManager(secret string, ttl time.Duration) *JWTManager {
	if secret == "" {
		secret = defaultJWTSecret
	}

	return &JWTManager{
		secret: []byte(secret),
		ttl:    ttl,
	}
}

func (m *JWTManager) GenerateToken(user *models.User) (string, error) {
	if user == nil || user.ID == 0 {
		return "", errors.New("user is required")
	}

	headerJSON, err := json.Marshal(map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	})
	if err != nil {
		return "", fmt.Errorf("marshal header: %w", err)
	}

	claimsJSON, err := json.Marshal(Claims{
		UserID:    user.ID,
		Email:     user.Email,
		ExpiresAt: time.Now().Add(m.ttl).Unix(),
	})
	if err != nil {
		return "", fmt.Errorf("marshal claims: %w", err)
	}

	header := base64.RawURLEncoding.EncodeToString(headerJSON)
	payload := base64.RawURLEncoding.EncodeToString(claimsJSON)
	signingInput := header + "." + payload
	signature := base64.RawURLEncoding.EncodeToString(m.sign(signingInput))

	return signingInput + "." + signature, nil
}

func (m *JWTManager) ParseToken(token string) (*Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, ErrInvalidToken
	}

	signingInput := parts[0] + "." + parts[1]
	expectedSignature := m.sign(signingInput)
	receivedSignature, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return nil, ErrInvalidToken
	}

	if !hmac.Equal(receivedSignature, expectedSignature) {
		return nil, ErrInvalidToken
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, ErrInvalidToken
	}

	var claims Claims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, ErrInvalidToken
	}

	if claims.UserID == 0 {
		return nil, ErrInvalidToken
	}

	if time.Now().Unix() >= claims.ExpiresAt {
		return nil, ErrExpiredToken
	}

	return &claims, nil
}

func (m *JWTManager) sign(input string) []byte {
	mac := hmac.New(sha256.New, m.secret)
	mac.Write([]byte(input))
	return mac.Sum(nil)
}
