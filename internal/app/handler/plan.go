package handler

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/novaru/billing-service/internal/app/service"
	E "github.com/novaru/billing-service/internal/shared/errors"
	"github.com/novaru/billing-service/internal/shared/response"
	"github.com/novaru/billing-service/pkg/logger"
)

type CreatePlanRequest struct {
	Slug        string         `json:"slug"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	PriceCents  int64          `json:"price_cents"`
	Currency    string         `json:"currency"`
	Interval    string         `json:"interval"`
	QuotaLimits map[string]any `json:"quota_limits"`
	Meta        map[string]any `json:"meta"`
}

type PlanHandler struct {
	service service.PlanService
}

func NewPlanHandler(s service.PlanService) *PlanHandler {
	return &PlanHandler{service: s}
}

func (h *PlanHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreatePlanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, E.NewInvalidInputError("invalid JSON format", err))
		return
	}

	plan, err := h.service.Create(
		r.Context(),
		req.Slug,
		req.Name,
		req.Description,
		req.PriceCents,
		req.Currency,
		req.Interval,
		req.QuotaLimits,
		req.Meta,
	)
	if err != nil {
		logger.Debug("could not create plan", zap.Error(err))
		response.WriteError(w, E.NewInvalidInputError("could not create plan", err))
		return
	}

	response.WriteCreated(w, plan)
}
