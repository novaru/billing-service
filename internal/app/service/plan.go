package service

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/novaru/billing-service/db/generated"
	"github.com/novaru/billing-service/internal/app/repository"
)

type PlanResponse struct {
	ID          string         `json:"id"`
	Slug        string         `json:"slug"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	PriceCents  int64          `json:"price_cents"`
	Currency    string         `json:"currency"`
	Interval    string         `json:"interval"`
	QuotaLimits map[string]any `json:"quota_limits"`
	Meta        map[string]any `json:"meta"`
}

type PlanService interface {
	Create(ctx context.Context, slug, name, description string, priceCents int64, currency, interval string, quotaLimits, meta map[string]any) (PlanResponse, error)
}

type planService struct {
	repo repository.PlanRepository
}

func NewPlanService(repo repository.PlanRepository) PlanService {
	return &planService{repo: repo}
}

func (s *planService) Create(ctx context.Context, slug, name, description string, priceCents int64, currency, interval string, quotaLimits, meta map[string]any) (PlanResponse, error) {
	quotaLimitsBytes, err := json.Marshal(quotaLimits)
	if err != nil {
		return PlanResponse{}, err
	}
	metaBytes, err := json.Marshal(meta)
	if err != nil {
		return PlanResponse{}, err
	}
	descriptionText := pgtype.Text{String: description, Valid: description != ""}

	plan, err := s.repo.Create(ctx, generated.CreatePlanParams{
		Slug:        slug,
		Name:        name,
		Description: descriptionText,
		PriceCents:  priceCents,
		Currency:    currency,
		Interval:    interval,
		QuotaLimits: quotaLimitsBytes,
		Meta:        metaBytes,
	})
	if err != nil {
		return PlanResponse{}, err
	}

	return PlanResponse{
		ID:          plan.ID.String(),
		Slug:        plan.Slug,
		Name:        plan.Name,
		Description: plan.Description.String,
		PriceCents:  plan.PriceCents,
		Currency:    plan.Currency,
		Interval:    plan.Interval,
		QuotaLimits: quotaLimits,
		Meta:        meta,
	}, nil
}
