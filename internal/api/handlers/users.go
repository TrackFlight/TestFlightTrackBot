package handlers

import (
	"encoding/json"
	"errors"
	"github.com/Laky-64/TestFlightTrackBot/internal/api/middleware"
	"github.com/Laky-64/TestFlightTrackBot/internal/api/types"
	"github.com/Laky-64/TestFlightTrackBot/internal/api/utils"
	"github.com/Laky-64/TestFlightTrackBot/internal/config"
	"github.com/Laky-64/TestFlightTrackBot/internal/db"
	"github.com/Laky-64/TestFlightTrackBot/internal/testflight"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"net/http"
	"strconv"
)

func GetLinks(dbCtx *db.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		list, err := dbCtx.ChatStore.TrackedList(
			r.Context().Value(middleware.UserIDKey).(int64),
		)
		if err != nil {
			utils.JSONError(w, types.ErrInternalServer, "Error fetching links", http.StatusInternalServerError)
			return
		}
		var result []types.Link
		for _, item := range list {
			var timestamp int64
			if item.LastAvailability.Valid {
				timestamp = item.LastAvailability.Time.UTC().Unix()
			}
			result = append(result, types.Link{
				ID:               item.ID,
				Tag:              utils.EncodeTag(item.ID),
				AppName:          item.AppName,
				IconURL:          item.IconURL,
				Description:      item.Description,
				Status:           string(item.Status),
				LastAvailability: timestamp,
			})
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(result)
	}
}

func DeleteLink(dbCtx *db.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		linkID, _ := strconv.Atoi(chi.URLParam(r, "linkID"))
		if linkID == 0 {
			utils.JSONError(w, types.ErrBadRequest, "Invalid link ID", http.StatusBadRequest)
			return
		}

		err := dbCtx.ChatStore.Delete(r.Context().Value(middleware.UserIDKey).(int64), int64(linkID))
		if errors.Is(err, pgx.ErrNoRows) {
			utils.JSONError(w, types.ErrBadRequest, "Link not found", http.StatusNotFound)
			return
		} else if err != nil {
			utils.JSONError(w, types.ErrInternalServer, "Error deleting link", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func AddLink(dbCtx *db.DB, cfg *config.Config) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var newLink types.NewLink
		if err := json.NewDecoder(r.Body).Decode(&newLink); err != nil {
			utils.JSONError(w, types.ErrBadRequest, "Invalid request body", http.StatusBadRequest)
			return
		}
		if newLink.Link == "" && newLink.ID == 0 {
			utils.JSONError(w, types.ErrBadRequest, "Link or ID must be provided", http.StatusBadRequest)
			return
		} else if newLink.Link != "" && newLink.ID != 0 {
			utils.JSONError(w, types.ErrBadRequest, "Only one of Link or ID should be provided", http.StatusBadRequest)
			return
		}

		if newLink.Link != "" {
			if !testflight.RegexLink.MatchString(newLink.Link) {
				utils.JSONError(w, types.ErrBadRequest, "Invalid TestFlight link format", http.StatusBadRequest)
				return
			}
		}

		if following, err := dbCtx.ChatStore.Track(r.Context().Value(middleware.UserIDKey).(int64), newLink.ID, newLink.Link, cfg.LimitFree); err != nil {
			utils.JSONError(w, types.ErrInternalServer, "Error tracking link", http.StatusInternalServerError)
			return
		} else if following == nil {
			utils.JSONError(w, types.ErrLinkAlreadyFollowing, "Link already tracked", http.StatusConflict)
			return
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			var timestamp int64
			if following.LastAvailability.Valid {
				timestamp = following.LastAvailability.Time.UTC().Unix()
			}
			_ = json.NewEncoder(w).Encode(types.Link{
				ID:               following.ID,
				AppID:            following.AppID.Int64,
				Tag:              utils.EncodeTag(following.ID),
				AppName:          following.AppName,
				IconURL:          following.IconURL,
				Description:      following.Description,
				Status:           string(following.Status),
				LastAvailability: timestamp,
			})
		}
	}
}
