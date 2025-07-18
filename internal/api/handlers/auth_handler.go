package handlers

import (
	"encoding/json"
	"github.com/Laky-64/TestFlightTrackBot/internal/api/types"
	"github.com/Laky-64/TestFlightTrackBot/internal/api/utils"
	"github.com/Laky-64/TestFlightTrackBot/internal/auth"
	"net/http"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func AuthHandler(botToken string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			utils.JSONError(w, types.ErrBadRequest, "Invalid request data", http.StatusBadRequest)
			return
		}
		initData := r.FormValue("initData")
		userID, valid := auth.ValidateInitData(initData, botToken)
		if !valid {
			utils.JSONError(w, types.ErrUnauthorized, "Invalid init data", http.StatusUnauthorized)
			return
		}
		token, err := auth.GenerateJWT(userID)
		if err != nil {
			utils.JSONError(w, types.ErrInternalServer, "Error generating token", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{
			"token": token,
		})
	}
}
