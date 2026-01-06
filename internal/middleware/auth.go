package middleware

import (
	"context"
	"net/http"

	"github.com/akramboussanni/marchive/config"
	"github.com/akramboussanni/marchive/internal/api"
	"github.com/akramboussanni/marchive/internal/jwt"
	"github.com/akramboussanni/marchive/internal/model"
	"github.com/akramboussanni/marchive/internal/repo"
	"github.com/akramboussanni/marchive/internal/utils"
	"github.com/go-chi/chi/v5"
)

func AddAuth(r chi.Router, ur *repo.UserRepo, tr *repo.TokenRepo) {
	r.Use(func(next http.Handler) http.Handler {
		return JWTAuth(config.JwtSecretBytes, ur, tr, model.CredentialJwt)(next)
	})
}

func AddOptionalAuth(r chi.Router, ur *repo.UserRepo, tr *repo.TokenRepo) {
	r.Use(func(next http.Handler) http.Handler {
		return OptionalJWTAuth(config.JwtSecretBytes, ur, tr, model.CredentialJwt)(next)
	})
}

func OptionalJWTAuth(secret []byte, ur *repo.UserRepo, tr *repo.TokenRepo, expectedType model.JwtType) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Try to get session cookie, but don't fail if not present
			sessionCookie, err := r.Cookie("session")
			if err != nil {
				// No session cookie, continue without user context
				next.ServeHTTP(w, r)
				return
			}

			// Try to validate token
			claims, err := jwt.ValidateToken(sessionCookie.Value, secret, tr)
			if err != nil {
				// Invalid token, continue without user context
				next.ServeHTTP(w, r)
				return
			}

			// Check token type
			if claims.Type != expectedType {
				// Wrong token type, continue without user context
				next.ServeHTTP(w, r)
				return
			}

			// Try to get user
			user, err := ur.GetUserByID(r.Context(), claims.UserID)
			if err != nil {
				// User not found, continue without user context
				next.ServeHTTP(w, r)
				return
			}

			// Check session ID
			if claims.SessionID != user.JwtSessionID {
				// Session mismatch, continue without user context
				next.ServeHTTP(w, r)
				return
			}

			// Valid user, add to context
			ctx := context.WithValue(r.Context(), utils.UserKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func JWTAuth(secret []byte, ur *repo.UserRepo, tr *repo.TokenRepo, expectedType model.JwtType) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := GetClaimsFromCookie(w, r, secret, ur, tr)
			if claims == nil {
				return
			}

			if claims.Type != expectedType {
				api.WriteInvalidCredentials(w)
				return
			}

			user, err := ur.GetUserByID(r.Context(), claims.UserID)
			if err != nil {
				api.WriteInternalError(w)
				return
			}

			if claims.SessionID != user.JwtSessionID {
				api.WriteInvalidCredentials(w)
				return
			}

			ctx := context.WithValue(r.Context(), utils.UserKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetClaimsFromCookie(w http.ResponseWriter, r *http.Request, secret []byte, ur *repo.UserRepo, tr *repo.TokenRepo) *jwt.Claims {
	sessionCookie, err := r.Cookie("session")
	if err != nil {
		api.WriteInvalidCredentials(w)
		return nil
	}

	return GetClaims(w, r, sessionCookie.Value, secret, tr)
}

func GetClaims(w http.ResponseWriter, r *http.Request, token string, secret []byte, tr *repo.TokenRepo) *jwt.Claims {
	claims, err := jwt.ValidateToken(token, config.JwtSecretBytes, tr)
	if err != nil {
		api.WriteInvalidCredentials(w)
		return nil
	}

	return claims
}

func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := utils.UserFromContext(r.Context())
		if !ok {
			api.WriteInvalidCredentials(w)
			return
		}

		if user.Role != "admin" {
			api.WriteMessage(w, http.StatusForbidden, "error", "admin access required")
			return
		}

		next.ServeHTTP(w, r)
	})
}
