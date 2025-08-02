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
	"github.com/jackc/pgx/v5"
	"maps"
	"net/http"
	"slices"
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
		var orderKeys []int64
		result := make(map[int64]*types.App)
		for _, item := range list {
			if _, exists := result[item.EntityID]; !exists {
				result[item.EntityID] = &types.App{
					ID: item.EntityID,
				}
				orderKeys = append(orderKeys, item.EntityID)
			}
			var timestamp int64
			if item.LastAvailability.Valid {
				timestamp = item.LastAvailability.Time.UTC().Unix()
			}
			result[item.EntityID].Links = append(result[item.EntityID].Links, types.Link{
				ID:               item.ID,
				URL:              item.LinkURL,
				Status:           item.Status,
				IsPublic:         item.IsPublic,
				LastAvailability: timestamp,
				LastUpdate:       item.LastUpdate.Time.UTC().Unix(),
			})
		}

		if len(result) > 0 {
			entityList, errAppList := dbCtx.AppStore.GetAppsInfo(
				slices.Collect(maps.Keys(result)),
			)
			if errAppList != nil {
				utils.JSONError(w, types.ErrInternalServer, "Error fetching apps info", http.StatusInternalServerError)
				return
			}
			for _, entityInfo := range entityList {
				entity := result[entityInfo.EntityID]
				entity.Name = entityInfo.AppName
				entity.IconURL = entityInfo.IconURL
				entity.Description = entityInfo.Description
				entity.Followers = entityInfo.Followers
			}
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
			entityList, errAppInfo := dbCtx.AppStore.GetAppsInfo(
				[]int64{following.EntityID},
			)
			if errAppInfo != nil || len(entityList) == 0 {
				utils.JSONError(w, types.ErrInternalServer, "Error fetching apps info", http.StatusInternalServerError)
				return
			}
			entityInfo := entityList[0]
			_ = json.NewEncoder(w).Encode(types.App{
				ID:          following.EntityID,
				Name:        entityInfo.AppName,
				IconURL:     entityInfo.IconURL,
				Description: entityInfo.Description,
				Followers:   entityInfo.Followers,
				Links: []types.Link{
					{
						ID:               following.ID,
						URL:              following.LinkURL,
						Status:           following.Status,
						IsPublic:         following.IsPublic,
						LastAvailability: timestamp,
						LastUpdate:       following.LastUpdate.Time.UTC().Unix(),
					},
				},
			})
		}
	}
}
