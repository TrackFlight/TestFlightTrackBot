package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/TrackFlight/TestFlightTrackBot/internal/api/types"
	"github.com/TrackFlight/TestFlightTrackBot/internal/api/utils"
	"github.com/TrackFlight/TestFlightTrackBot/internal/db"
)

func SearchApps(dbCtx *db.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "" {
			utils.JSONError(w, types.ErrBadRequest, "App name is required", http.StatusBadRequest)
			return
		}
		appIDs, err := dbCtx.AppStore.Search(name)
		if err != nil {
			utils.JSONError(w, types.ErrInternalServer, "Error fetching apps", http.StatusInternalServerError)
			return
		}
		result := make(map[int64]*types.TrendingApp)
		if len(appIDs) > 0 {
			appsInfo, err := dbCtx.AppStore.GetAppsInfo(appIDs)
			if err != nil {
				utils.JSONError(w, types.ErrInternalServer, "Error fetching app info", http.StatusInternalServerError)
				return
			}
			for _, app := range appsInfo {
				result[app.EntityID] = &types.TrendingApp{
					BaseApp: types.BaseApp{
						ID:          app.EntityID,
						Name:        app.AppName,
						IconURL:     app.IconURL,
						Description: app.Description,
						Followers:   app.Followers,
					},
				}
			}
			appLinks, err := dbCtx.LinkStore.GetLinksByApps(appIDs)
			if err != nil {
				utils.JSONError(w, types.ErrInternalServer, "Error fetching app links", http.StatusInternalServerError)
				return
			}
			for _, link := range appLinks {
				if app, exists := result[link.AppID]; exists {
					var timestamp int64
					if link.LastAvailability.Valid {
						timestamp = link.LastAvailability.Time.UTC().Unix()
					}
					app.Links = append(app.Links, types.BaseLink{
						ID:               link.ID,
						URL:              link.URL,
						Status:           link.Status,
						IsPublic:         link.IsPublic,
						LastAvailability: timestamp,
						LastUpdate:       link.LastUpdate.Time.UTC().Unix(),
					})
				}
			}
		}

		orderedResult := make([]*types.TrendingApp, 0, len(appIDs))
		for _, key := range appIDs {
			if app, exists := result[key]; exists {
				orderedResult = append(orderedResult, app)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(orderedResult)
	}
}

func GetTrendingApps(dbCtx *db.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		trendingAppIDs, err := dbCtx.AppStore.GetTrending()
		if err != nil {
			utils.JSONError(w, types.ErrInternalServer, "Error fetching trending apps", http.StatusInternalServerError)
			return
		}

		appsInfo, err := dbCtx.AppStore.GetAppsInfo(trendingAppIDs)
		if err != nil {
			utils.JSONError(w, types.ErrInternalServer, "Error fetching app info", http.StatusInternalServerError)
			return
		}

		result := make(map[int64]*types.TrendingApp)
		for _, app := range appsInfo {
			result[app.EntityID] = &types.TrendingApp{
				BaseApp: types.BaseApp{
					ID:          app.EntityID,
					Name:        app.AppName,
					IconURL:     app.IconURL,
					Description: app.Description,
					Followers:   app.Followers,
				},
			}
		}

		trendingAppsLinks, err := dbCtx.LinkStore.GetLinksByApps(trendingAppIDs)
		if err != nil {
			utils.JSONError(w, types.ErrInternalServer, "Error fetching links for trending apps", http.StatusInternalServerError)
			return
		}
		for _, link := range trendingAppsLinks {
			if app, exists := result[link.AppID]; exists {
				var timestamp int64
				if link.LastAvailability.Valid {
					timestamp = link.LastAvailability.Time.UTC().Unix()
				}
				app.Links = append(app.Links, types.BaseLink{
					ID:               link.ID,
					URL:              link.URL,
					Status:           link.Status,
					IsPublic:         link.IsPublic,
					LastAvailability: timestamp,
					LastUpdate:       link.LastUpdate.Time.UTC().Unix(),
				})
			}
		}

		var orderedResult []*types.TrendingApp
		for _, key := range trendingAppIDs {
			if app, exists := result[key]; exists {
				orderedResult = append(orderedResult, app)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(orderedResult)
	}
}
