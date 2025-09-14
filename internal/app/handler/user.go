package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/novaru/billing-service/internal/app/service"
	E "github.com/novaru/billing-service/internal/shared/errors"
	"github.com/novaru/billing-service/internal/shared/response"
	"github.com/novaru/billing-service/pkg/logger"
)

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(s service.UserService) *UserHandler {
	return &UserHandler{service: s}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		response.WriteError(w, E.NewInvalidInputError("invalid JSON format", err))
		return
	}

	user, err := h.service.Create(r.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	response.WriteCreated(w, user)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, E.NewInvalidInputError("invalid request body", err))
		return
	}

	token, expiresAt, err := h.service.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		response.WriteError(w, E.NewUnauthorizedError("login failed", err))
		return
	}

	response.WriteSuccess(w, map[string]any{
		"token":      token,
		"expires_at": expiresAt.Format(time.RFC3339),
		"token_type": "Bearer",
	})
}

func (h *UserHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 10
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		offset = 0
	}

	users, err := h.service.FindAll(r.Context(), int32(limit), int32(offset))
	if err != nil {
		response.WriteError(w, err)
		return
	}

	response.WriteSuccess(w, users)
}

func (h *UserHandler) FindByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	userID, err := uuid.Parse(idParam)
	if err != nil {
		http.Error(w, "invalid UUID", http.StatusBadRequest)
		return
	}

	user, err := h.service.FindByID(r.Context(), userID)
	if err != nil {
		if errors.Is(err, E.ErrNotFound) {
			response.WriteError(w, E.NewNotFoundError("user", "user with given ID does not exist"))
			return
		}

		logger.Debug("failed to fetch user:", zap.Error(err))
		response.WriteError(w, E.ErrInternal)
		return
	}

	response.WriteSuccess(w, user)
}
