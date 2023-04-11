package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/vitalis-virtus/rest-api-metamask/storage"
	"github.com/vitalis-virtus/rest-api-metamask/utils"
	"github.com/vitalis-virtus/rest-api-metamask/utils/fail"
)

func AuthMiddleware(storage *storage.MemoryStorage, jwtProvider *utils.JwtHmacProvider) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			headerValue := r.Header.Get("Authorization")
			const prefix = "Bearer "
			if len(headerValue) < len(prefix) {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			tokenString := headerValue[len(prefix):]
			if len(tokenString) == 0 {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			claims, err := jwtProvider.Verify(tokenString)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			user, err := storage.Get(claims.Subject)
			if err != nil {
				if errors.Is(err, fail.ErrUserNotExists) {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			ctx := context.WithValue(r.Context(), "user", user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
