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

type PlanRepository interface {
	Create(ctx context.Context, arg generated.CreatePlanParams) (generated.Plan, error)
	FindAll(ctx context.Context) ([]generated.Plan, error)
	FindBySlug(ctx context.Context, slug string) (generated.Plan, error)
}

type planRepository struct {
	q *generated.Queries
}

func NewPlanRepository(q *generated.Queries) PlanRepository {
	return &planRepository{q: q}
}

func (r *planRepository) Create(ctx context.Context, arg generated.CreatePlanParams) (generated.Plan, error) {
	id, err := uuid.NewV7()
	if err != nil {
		logger.Fatal("failed to generate uuid:", zap.Error(err))
	}
	return r.q.CreatePlan(ctx, generated.CreatePlanParams{
		ID:          id,
		Slug:        arg.Slug,
		Name:        arg.Name,
		Description: arg.Description,
		PriceCents:  arg.PriceCents,
		Currency:    arg.Currency,
		Interval:    arg.Interval,
		QuotaLimits: arg.QuotaLimits,
		Meta:        arg.Meta,
	})
}

func (r *planRepository) FindAll(ctx context.Context) ([]generated.Plan, error) {
	return r.q.ListPlans(ctx)
}

func (r *planRepository) FindBySlug(ctx context.Context, slug string) (generated.Plan, error) {
	logger.Debug("retrieving plan by slug", zap.String("slug", slug))

	plan, err := r.q.GetPlanBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Debug("plan not found", zap.String("slug", slug))
			return generated.Plan{}, E.ErrNotFound
		}

		logger.Error("failed to retrieve plan by slug", zap.String("slug", slug), zap.Error(err))
		return generated.Plan{}, err
	}

	logger.Debug("plan retrieved successfully", zap.String("slug", slug))
	return plan, nil
}
