package handlers

import (
	"encoding/json"
	"github.com/Laky-64/TestFlightTrackBot/internal/api/types"
	"github.com/Laky-64/TestFlightTrackBot/internal/db"
	"github.com/jackc/pgx/v5/pgtype"
	"net/http"
)

func GetTrendingApps(dbCtx *db.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		trendingApps, err := dbCtx.AppStore.GetTrending()
		if err != nil {
			http.Error(w, "Error fetching trending apps", http.StatusInternalServerError)
			return
		}

		var orderKeys []int64
		result := make(map[int64]*types.TrendingApp)
		for _, app := range trendingApps {
			result[app.ID] = &types.TrendingApp{
				BaseApp: types.BaseApp{
					ID: app.ID,
					Name: pgtype.Text{
						String: app.AppName,
						Valid:  true,
					},
					IconURL:     app.IconURL,
					Description: app.Description,
					Followers:   app.Followers,
				},
			}
			orderKeys = append(orderKeys, app.ID)
		}

		trendingAppsLinks, err := dbCtx.LinkStore.GetLinksByApps(orderKeys)
		if err != nil {
			http.Error(w, "Error fetching links for trending apps", http.StatusInternalServerError)
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
