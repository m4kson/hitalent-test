package service

import (
	"fmt"
	"sync"
	"time"

	"hitalent-test/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

type TokenClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type TokenService struct {
	cfg *config.JWTConfig
}

type RefreshTokenInfo struct {
	UserID    string
	ExpiresAt time.Time
}

type RefreshTokenStore struct {
	tokens map[string]*RefreshTokenInfo
	mu     sync.RWMutex
}

func NewTokenService(cfg *config.JWTConfig) *TokenService {
	return &TokenService{cfg: cfg}
}

func NewRefreshTokenStore() *RefreshTokenStore {
	return &RefreshTokenStore{
		tokens: make(map[string]*RefreshTokenInfo),
	}
}

func (s *TokenService) GenerateAccessToken(userID, email string) (string, error) {
	claims := TokenClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.cfg.AccessTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.Secret))
}

func (s *TokenService) GenerateRefreshToken(userID string) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.cfg.RefreshTokenExpiry)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.Secret))
}

func (s *TokenService) VerifyToken(tokenString string) (*TokenClaims, error) {
	claims := &TokenClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.cfg.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func (store *RefreshTokenStore) Save(token string, info *RefreshTokenInfo) {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.tokens[token] = info
}

func (store *RefreshTokenStore) Get(token string) (*RefreshTokenInfo, bool) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	info, ok := store.tokens[token]
	if !ok {
		return nil, false
	}

	if info.ExpiresAt.Before(time.Now()) {
		return nil, false
	}
	return info, true
}

func (store *RefreshTokenStore) Delete(token string) {
	store.mu.Lock()
	defer store.mu.Unlock()
	delete(store.tokens, token)
}

func (store *RefreshTokenStore) CleanupExpired() {
	store.mu.Lock()
	defer store.mu.Unlock()

	now := time.Now()
	for token, info := range store.tokens {
		if info.ExpiresAt.Before(now) {
			delete(store.tokens, token)
		}
	}
}
