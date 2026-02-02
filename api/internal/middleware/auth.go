package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	"github.com/yourusername/hireme/api/internal/config"
	"github.com/yourusername/hireme/api/pkg/httputil"
)

type contextKey string

const (
	// UserIDKey is the context key for the authenticated user ID
	UserIDKey contextKey = "userID"
)

// Auth creates authentication middleware
func Auth(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Development bypass
			if cfg.Auth.BypassEnabled {
				slog.Debug("auth bypass enabled",
					"user_id", cfg.Auth.BypassUserID,
				)
				ctx := context.WithValue(r.Context(), UserIDKey, cfg.Auth.BypassUserID)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			// Extract token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				httputil.Error(w, http.StatusUnauthorized, "missing authorization header")
				return
			}

			// Expect "Bearer <token>"
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
				httputil.Error(w, http.StatusUnauthorized, "invalid authorization header format")
				return
			}

			tokenString := parts[1]

			// Parse and validate JWT
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				// Validate signing method
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(cfg.Auth.JWTSecret), nil
			})

			if err != nil {
				slog.Debug("jwt validation failed", "error", err)
				httputil.Error(w, http.StatusUnauthorized, "invalid token")
				return
			}

			if !token.Valid {
				httputil.Error(w, http.StatusUnauthorized, "invalid token")
				return
			}

			// Extract user ID from claims
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				httputil.Error(w, http.StatusUnauthorized, "invalid token claims")
				return
			}

			userID, ok := claims["sub"].(string)
			if !ok || userID == "" {
				httputil.Error(w, http.StatusUnauthorized, "missing user ID in token")
				return
			}

			// Add user ID to context
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserID extracts the user ID from the request context
func GetUserID(ctx context.Context) string {
	if userID, ok := ctx.Value(UserIDKey).(string); ok {
		return userID
	}
	return ""
}

// MustGetUserID extracts the user ID from context, panics if not found
// Use only when auth middleware has already validated the request
func MustGetUserID(ctx context.Context) string {
	userID := GetUserID(ctx)
	if userID == "" {
		panic("user ID not found in context - auth middleware not applied?")
	}
	return userID
}
