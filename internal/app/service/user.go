package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/novaru/billing-service/db/generated"
	"github.com/novaru/billing-service/internal/app/repository"
	"github.com/novaru/billing-service/internal/config"
	E "github.com/novaru/billing-service/internal/shared/errors"
)

type UserResponse struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

type UserService interface {
	Login(ctx context.Context, email, password string) (string, time.Time, error)
	RefreshToken(ctx context.Context, oldToken string) (string, time.Time, error)
	Create(ctx context.Context, name, email, password string) (UserResponse, error)
	FindAll(ctx context.Context, limit, offset int32) ([]UserResponse, error)
	FindByID(ctx context.Context, id uuid.UUID) (UserResponse, error)
}

type userService struct {
	cfg  *config.Config
	repo repository.UserRepository
}

func NewUserService(cfg *config.Config, repo repository.UserRepository) UserService {
	return &userService{cfg: cfg, repo: repo}
}

func (s *userService) Login(ctx context.Context, email, password string) (string, time.Time, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return "", time.Time{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", time.Time{}, errors.New("invalid credentials")
	}

	return s.generateNewToken(user.ID.String(), user.Email)
}

func (s *userService) RefreshToken(ctx context.Context, oldToken string) (string, time.Time, error) {
	token, err := jwt.Parse(oldToken, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(s.cfg.JWTSecret), nil
	})

	if err != nil {
		return "", time.Time{}, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", time.Time{}, errors.New("invalid claims")
	}

	// Extract claims
	userID, ok := claims["sub"].(string)
	if !ok {
		return "", time.Time{}, errors.New("invalid user ID in token")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return "", time.Time{}, errors.New("invalid email in token")
	}

	// Check if token is close to expiry (within 1 hour) or already expired
	exp := int64(claims["exp"].(float64))
	now := time.Now().Unix()

	// Allow refresh if token expires within 1 hour or is already expired (but not too old)
	if now < exp-3600 {
		return "", time.Time{}, errors.New("token not eligible for refresh yet")
	}

	// Don't allow refresh of very old tokens (older than 7 days)
	if now > exp+604800 {
		return "", time.Time{}, errors.New("token too old to refresh")
	}

	return s.generateNewToken(userID, email)
}

func (s *userService) Create(ctx context.Context, name, email, password string) (UserResponse, error) {
	if err := s.Validate(name, email, password); err != nil {
		return UserResponse{}, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return UserResponse{}, err
	}

	user, err := s.repo.Create(ctx, name, email, string(hash))
	if err != nil {
		return UserResponse{}, err
	}

	return UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (s *userService) FindAll(ctx context.Context, limit, offset int32) ([]UserResponse, error) {
	users, err := s.repo.FindAll(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	var userResponses []UserResponse
	for _, user := range users {
		userResponses = append(userResponses, UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		})
	}

	return userResponses, nil
}

func (s *userService) FindByID(ctx context.Context, id uuid.UUID) (UserResponse, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return UserResponse{}, err
	}
	return s.convertToResponse(user), nil
}

func (s *userService) Validate(name, email, password string) error {
	if strings.TrimSpace(name) == "" {
		return E.NewInvalidInputError("name is required", nil)
	}
	if strings.TrimSpace(email) == "" {
		return E.NewInvalidInputError("email is required", nil)
	}
	if len(password) < 6 {
		return E.NewInvalidInputError("password must be at least 3 characters", nil)
	}
	return nil
}

func (s *userService) convertToResponse(user generated.User) UserResponse {
	return UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
}

func (s *userService) generateNewToken(userID string, email string) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(24 * time.Hour)

	claims := jwt.MapClaims{
		"sub":   userID,
		"email": email,
		"exp":   expiresAt.Unix(),
		"iat":   now.Unix(),
		"iss":   "billing-service",
		"aud":   "billing-api",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, expiresAt, nil
}
