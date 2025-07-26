package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Laky-64/TestFlightTrackBot/internal/api/middleware"
	"github.com/Laky-64/TestFlightTrackBot/internal/api/types"
	"github.com/Laky-64/TestFlightTrackBot/internal/api/utils"
	"github.com/Laky-64/TestFlightTrackBot/internal/config"
	"github.com/Laky-64/TestFlightTrackBot/internal/db"
	"github.com/Laky-64/TestFlightTrackBot/internal/testflight"
	"github.com/jackc/pgx/v5"
	"net/http"
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

		noAppCounter := 0
		var orderKeys []string
		result := make(map[string]*types.App)
		for _, item := range list {
			var key string
			if item.AppID.Valid {
				key = fmt.Sprintf("app_%d", item.AppID.Int64)
			} else {
				noAppCounter++
				key = fmt.Sprintf("no_app_%d", noAppCounter)
			}

			if _, exists := result[key]; !exists {
				result[key] = &types.App{
					ID:          item.AppID,
					Name:        item.AppName,
					IconURL:     item.IconURL,
					Description: item.Description,
				}
				orderKeys = append(orderKeys, key)
			}
			var timestamp int64
			if item.LastAvailability.Valid {
				timestamp = item.LastAvailability.Time.UTC().Unix()
			}
			result[key].Links = append(result[key].Links, types.Link{
				ID:               item.ID,
				URL:              item.LinkURL,
				Status:           item.Status,
				LastAvailability: timestamp,
				LastUpdate:       item.LastUpdate.Time.UTC().Unix(),
			})
		}

		var orderedResult []*types.App
		for _, key := range orderKeys {
			if app, exists := result[key]; exists {
				orderedResult = append(orderedResult, app)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(orderedResult)
	}
}

func DeleteLinks(dbCtx *db.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var removeLinks []int64
		if err := json.NewDecoder(r.Body).Decode(&removeLinks); err != nil {
			utils.JSONError(w, types.ErrBadRequest, "Invalid request body", http.StatusBadRequest)
			return
		}
		if len(removeLinks) == 0 {
			utils.JSONError(w, types.ErrBadRequest, "No links provided", http.StatusBadRequest)
			return
		}
		removeLinks = removeLinks[:min(len(removeLinks), 5)]

		err := dbCtx.ChatStore.Delete(r.Context().Value(middleware.UserIDKey).(int64), removeLinks)
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
		var newLink types.AddLinkRequest
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
			_ = json.NewEncoder(w).Encode(types.App{
				ID:          following.AppID,
				Name:        following.AppName,
				IconURL:     following.IconURL,
				Description: following.Description,
				Links: []types.Link{
					{
						ID:               following.ID,
						URL:              following.LinkURL,
						Status:           following.Status,
						LastAvailability: timestamp,
						LastUpdate:       following.LastUpdate.Time.UTC().Unix(),
					},
				},
			})
		}
	}
}
