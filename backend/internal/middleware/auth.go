package middleware

import (
	"context"
	"net/http"

	"github.com/akramboussanni/gocode/config"
	"github.com/akramboussanni/gocode/internal/api"
	"github.com/akramboussanni/gocode/internal/jwt"
	"github.com/akramboussanni/gocode/internal/model"
	"github.com/akramboussanni/gocode/internal/repo"
	"github.com/akramboussanni/gocode/internal/utils"
	"github.com/go-chi/chi/v5"
)

func AddAuth(r chi.Router, ur *repo.UserRepo, tr *repo.TokenRepo) {
	r.Use(func(next http.Handler) http.Handler {
		return JWTAuth(config.JwtSecretBytes, ur, tr, model.CredentialJwt)(next)
	})
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
