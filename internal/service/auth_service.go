package service

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"hitalent-test/internal/domain"
	"hitalent-test/internal/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo      repository.UserRepository
	tokenService  *TokenService
	refreshTokens *RefreshTokenStore
}

func NewAuthService(
	userRepo repository.UserRepository,
	tokenService *TokenService,
	refreshTokens *RefreshTokenStore,
) *AuthService {
	return &AuthService{
		userRepo:      userRepo,
		tokenService:  tokenService,
		refreshTokens: refreshTokens,
	}
}

func (s *AuthService) Register(email, password string) (*domain.User, error) {
	if err := validateEmail(email); err != nil {
		return nil, fmt.Errorf("%w: invalid email format", domain.ErrInvalidInput)
	}

	if len(password) < 8 {
		return nil, fmt.Errorf("%w: password must be at least 8 characters", domain.ErrInvalidInput)
	}

	_, err := s.userRepo.GetByEmail(email)
	if err == nil {
		return nil, fmt.Errorf("%w: user with this email already exists", domain.ErrInvalidInput)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &domain.User{
		ID:           uuid.New().String(),
		Email:        strings.ToLower(email),
		PasswordHash: string(hash),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *AuthService) Login(email, password string) (*domain.AuthResponse, error) {
	user, err := s.userRepo.GetByEmail(strings.ToLower(email))
	if err != nil {
		return nil, fmt.Errorf("%w: user not found", domain.ErrInvalidInput)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("%w: invalid credentials", domain.ErrInvalidInput)
	}

	accessToken, err := s.tokenService.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.tokenService.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	s.refreshTokens.Save(refreshToken, &RefreshTokenInfo{
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(s.tokenService.cfg.RefreshTokenExpiry),
	})

	return &domain.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

func (s *AuthService) RefreshAccessToken(refreshToken string) (string, error) {
	info, ok := s.refreshTokens.Get(refreshToken)
	if !ok {
		return "", fmt.Errorf("%w: invalid or expired refresh token", domain.ErrInvalidInput)
	}

	user, err := s.userRepo.GetByID(info.UserID)
	if err != nil {
		return "", fmt.Errorf("failed to get user: %w", err)
	}

	accessToken, err := s.tokenService.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return "", fmt.Errorf("failed to generate access token: %w", err)
	}

	return accessToken, nil
}

func validateEmail(email string) error {
	const emailRegex = `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	if !re.MatchString(email) {
		return fmt.Errorf("invalid email format")
	}
	return nil
}
