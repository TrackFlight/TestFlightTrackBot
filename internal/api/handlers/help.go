package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/TrackFlight/TestFlightTrackBot/internal/api/middleware"
	"github.com/TrackFlight/TestFlightTrackBot/internal/api/types"
	"github.com/TrackFlight/TestFlightTrackBot/internal/api/utils"
	"github.com/TrackFlight/TestFlightTrackBot/internal/config"
	"github.com/TrackFlight/TestFlightTrackBot/internal/db"
)

func GetConfig(cfg *config.Config) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(
			types.Config{
				LimitFree:         cfg.LimitFree,
				LimitPremium:      cfg.LimitPremium,
				MaxFollowingLinks: cfg.MaxFollowingLinks,
			},
		)
	}
}

func GetSettingsPreferences(dbCtx *db.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		preferences, err := dbCtx.PreferencesStore.GetNotificationsPreferences(
			r.Context().Value(middleware.UserIDKey).(int64),
		)
		if err != nil {
			utils.JSONError(w, types.ErrInternalServer, "Error fetching preferences", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(
			preferences,
		)
	}
}

func EditSettingsPreferences(dbCtx *db.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var req db.GetNotificationsPreferencesPreferencesRow
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.JSONError(w, types.ErrBadRequest, "Invalid request data", http.StatusBadRequest)
			return
		}
		err := dbCtx.PreferencesStore.UpdateNotificationsPreferences(
			r.Context().Value(middleware.UserIDKey).(int64),
			req.NewFeaturesNotifications,
			req.WeeklyInsightsNotifications,
		)
		if err != nil {
			utils.JSONError(w, types.ErrInternalServer, "Error updating preferences", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	}
}
