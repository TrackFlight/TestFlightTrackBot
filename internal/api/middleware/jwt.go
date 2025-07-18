package middleware

import (
	"context"
	"github.com/Laky-64/TestFlightTrackBot/internal/api/types"
	"github.com/Laky-64/TestFlightTrackBot/internal/api/utils"
	"github.com/Laky-64/TestFlightTrackBot/internal/auth"
	"net/http"
	"strconv"
	"strings"
)

const UserIDKey string = "userID"

func JWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			utils.JSONError(w, types.ErrUnauthorized, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := auth.ValidateJWT(tokenStr)
		if err != nil || !token.Valid {
			utils.JSONError(w, types.ErrUnauthorized, "Invalid or expired token", http.StatusUnauthorized)
			return
		}
		subject, err := token.Claims.GetSubject()
		if err != nil {
			utils.JSONError(w, types.ErrInternalServer, "Invalid token claims", http.StatusInternalServerError)
			return
		}

		userID, _ := strconv.Atoi(subject)
		ctx := context.WithValue(r.Context(), UserIDKey, int64(userID))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
