package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/ilhamgepe/simplebank/token"
	"github.com/ilhamgepe/simplebank/utils"
)

type contextKey string

const (
	authPayload contextKey = "authPayload"
)

func authMiddleware(config utils.Config, maker token.Maker) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authorizationHeader := r.Header.Get("Authorization")
			if authorizationHeader == "" {
				writeJSON(w, http.StatusUnauthorized, Response{
					Status:  false,
					Message: "unauthorized",
				})
				return
			}
			fields := strings.Fields(authorizationHeader)
			if len(fields) < 2 {
				writeJSON(w, http.StatusUnauthorized, Response{
					Status:  false,
					Message: "unauthorized",
				})
				return
			}
			if strings.ToLower(fields[0]) != "bearer" {
				writeJSON(w, http.StatusUnauthorized, Response{
					Status:  false,
					Message: "unsupported authorization header",
				})
				return
			}

			claim, err := maker.VerifyToken(fields[1])
			if err != nil {
				writeJSON(w, http.StatusUnauthorized, Response{
					Status:  false,
					Message: "unauthorized",
				})
				return
			}

			ctx := context.WithValue(r.Context(), authPayload, claim)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
