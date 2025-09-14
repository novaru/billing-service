package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/novaru/billing-service/internal/config"
	E "github.com/novaru/billing-service/internal/shared/errors"
	"github.com/novaru/billing-service/internal/shared/response"
)

type contextKey string

const userCtxKey contextKey = "user_id"

// AuthMiddleware validates JWT token and attaches user ID to request context.
func AuthMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.WriteError(w, E.NewUnauthorizedError("missing authorization header", nil))
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				response.WriteError(w, E.NewUnauthorizedError("invalid authorization header", nil))
				return
			}

			tokenStr := parts[1]

			// Parse and validate JWT
			token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
				// Ensure token uses expected signing method
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(cfg.JWTSecret), nil
			})

			if err != nil || !token.Valid {
				response.WriteError(w, E.NewUnauthorizedError("invalid or expired token", nil))
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				response.WriteError(w, E.NewUnauthorizedError("invalid token claims", nil))
				return
			}

			userID, ok := claims["sub"].(string)
			if !ok {
				response.WriteError(w, E.NewUnauthorizedError("invalid token subject", nil))
				return
			}

			// Put user ID into request context
			ctx := context.WithValue(r.Context(), userCtxKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func APIKeyAuth() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKey := r.Header.Get("Authorization")
			if apiKey == "" {
				http.Error(w, "API key required", http.StatusUnauthorized)
				return
			}

			// Validate API key, get customer info
			customer, err := validateAPIKey(apiKey)
			if err != nil {
				http.Error(w, "Invalid API key", http.StatusUnauthorized)
				return
			}

			// Add customer to context
			ctx := context.WithValue(r.Context(), "customer", customer)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserID extracts user ID from request context
func GetUserID(r *http.Request) (string, bool) {
	id, ok := r.Context().Value(userCtxKey).(string)
	return id, ok
}

func validateAPIKey(key string) (any, error) {
	// TODO: implement validation
	return key, nil
}
