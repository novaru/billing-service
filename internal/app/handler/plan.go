package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"go.uber.org/zap"

	"github.com/go-chi/chi/v5"
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

func (h *PlanHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	plans, err := h.service.FindAll(r.Context())
	if err != nil {
		logger.Debug("could not retrieve plans", zap.Error(err))
		response.WriteError(w, E.NewInternalError("could not retrieve plans", err))
		return
	}

	response.WriteSuccess(w, plans)
}

func (h *PlanHandler) FindBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		response.WriteError(w, E.NewInvalidInputError("slug is required", nil))
		return
	}

	plan, err := h.service.FindBySlug(r.Context(), slug)
	if err != nil {
		if errors.Is(err, E.ErrNotFound) {
			response.WriteError(w, E.NewNotFoundError("plan not found", "plan with given slug is not exist"))
			return
		}
		logger.Debug("could not retrieve plan", zap.String("slug", slug), zap.Error(err))
		response.WriteError(w, E.NewInternalError("could not retrieve plan", err))
		return
	}

	response.WriteSuccess(w, plan)
}
