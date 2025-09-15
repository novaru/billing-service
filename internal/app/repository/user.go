package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/novaru/billing-service/db/generated"
	E "github.com/novaru/billing-service/internal/shared/errors"
	"github.com/novaru/billing-service/pkg/logger"
)

type UserRepository interface {
	Create(ctx context.Context, name, email, password_hash string) (generated.User, error)
	FindAll(ctx context.Context, limit, offset int32) ([]generated.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (generated.User, error)
	FindByEmail(ctx context.Context, email string) (generated.User, error)
}

type userRepository struct {
	q *generated.Queries
}

func NewUserRepository(q *generated.Queries) UserRepository {
	return &userRepository{q: q}
}

func (r *userRepository) Create(ctx context.Context, name, email, password_hash string) (generated.User, error) {
	id, err := uuid.NewV7()
	if err != nil {
		logger.Fatal("failed to generate uuid:", zap.Error(err))
	}

	return r.q.CreateUser(ctx, generated.CreateUserParams{
		ID:           id,
		Name:         name,
		Email:        email,
		PasswordHash: password_hash,
	})
}

func (r *userRepository) FindAll(ctx context.Context, limit, offset int32) ([]generated.User, error) {
	return r.q.ListUsers(ctx, generated.ListUsersParams{
		Limit:  limit,
		Offset: offset,
	})
}

func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (generated.User, error) {
	logger.Debug("retrieving user by ID", zap.String("user_id", id.String()))

	user, err := r.q.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Debug("user not found", zap.String("user_id", id.String()))
			return generated.User{}, E.ErrNotFound
		}

		logger.Info("failed to get user by ID",
			zap.String("user_id", id.String()),
			zap.Error(err))
		return generated.User{}, err
	}

	logger.Debug("successfully retrieved user", zap.String("user_id", id.String()))
	return user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (generated.User, error) {
	logger.Debug("retrieving user by email", zap.String("email", email))

	user, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Debug("user not found", zap.String("email", email))
			return generated.User{}, E.ErrNotFound
		}

		logger.Info("failed to get user by email",
			zap.String("email", email),
			zap.Error(err))
		return generated.User{}, err
	}

	logger.Debug("successfully retrieved user", zap.String("email", email))
	return user, nil
}
