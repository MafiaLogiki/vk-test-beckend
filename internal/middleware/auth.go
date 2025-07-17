package middleware

import (
	"context"
	"net/http"
	"strconv"

	"marketplace-service/internal/token"
)

func isAuthorized(tok *token.Service, r *http.Request) bool {
	tokenStr := token.ExtractToken(r)
	_, err := tok.ValidateToken(tokenStr)
	return err == nil
}

func AuthMiddleware(tok *token.Service, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !isAuthorized(tok, r) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenStr := token.ExtractToken(r)
		tokenValid, _ := tok.ValidateToken(tokenStr)
		token, _ := strconv.Atoi(tokenValid)

		ctx := r.Context()
		ctx = context.WithValue(ctx, "token", token)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func OptionalAuthMiddleware(tok *token.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if token == "" {
				next.ServeHTTP(w, r)
				return
			}

			if !isAuthorized(tok, r) {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
