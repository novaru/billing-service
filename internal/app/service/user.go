package service

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/novaru/billing-service/db/generated"
	"github.com/novaru/billing-service/internal/app/repository"
	"github.com/novaru/billing-service/pkg/logger"
)

type UserResponse struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

type UserService interface {
	Create(ctx context.Context, name, email, password string) (UserResponse, error)
	FindAll(ctx context.Context, limit, offset int32) ([]UserResponse, error)
	FindByID(ctx context.Context, id uuid.UUID) (UserResponse, error)
	FindByEmail(ctx context.Context, email int) (UserResponse, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Create(ctx context.Context, name, email, password string) (UserResponse, error) {
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
	logger.Debug("retrieving user by ID", zap.String("user_id", id.String()))

	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		// Error already logged and wrapped in repository
		return UserResponse{}, err
	}

	logger.Debug("successfully retrieved user", zap.String("user_id", id.String()))

	return s.convertToResponse(user), nil
}

func (s *userService) FindByEmail(ctx context.Context, email int) (UserResponse, error) {
	// TODO: implement this method
	return UserResponse{}, nil
}

func (s *userService) convertToResponse(user generated.User) UserResponse {
	return UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
}
