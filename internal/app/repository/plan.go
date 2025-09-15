package repository

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/novaru/billing-service/db/generated"
	"github.com/novaru/billing-service/pkg/logger"
)

type PlanRepository interface {
	Create(ctx context.Context, arg generated.CreatePlanParams) (generated.Plan, error)
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
