package handlers

import (
	"encoding/json"
	"github.com/Laky-64/TestFlightTrackBot/internal/api/types"
	"github.com/Laky-64/TestFlightTrackBot/internal/config"
	"net/http"
)

func GetConfig(cfg *config.Config) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(
			types.Config{
				LimitFree:           cfg.LimitFree,
				LimitPremium:        cfg.LimitPremium,
				MaxFollowingPerUser: cfg.MaxFollowingPerUser,
			},
		)
	}
}
