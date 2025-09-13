package repository

import (
	"context"
	"log"

	"github.com/google/uuid"

	"github.com/novaru/billing-service/db/generated"
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
		log.Fatal("failed to generate uuid:", err)
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
	player, err := r.q.GetUserByID(ctx, id)
	if err != nil {
		return generated.User{}, err
	}

	return player, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (generated.User, error) {
	// TODO: implement find by email
	return generated.User{}, nil
}
