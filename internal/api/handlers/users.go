package handlers

import (
	"encoding/json"
	"errors"
	"github.com/Laky-64/TestFlightTrackBot/internal/api/middleware"
	"github.com/Laky-64/TestFlightTrackBot/internal/api/types"
	"github.com/Laky-64/TestFlightTrackBot/internal/api/utils"
	"github.com/Laky-64/TestFlightTrackBot/internal/db"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func GetLinks(dbCtx *db.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		list, err := dbCtx.ChatLinkStore.TrackedList(
			r.Context().Value(middleware.UserIDKey).(int64),
		)
		if err != nil {
			utils.JSONError(w, types.ErrInternalServer, "Error fetching links", http.StatusInternalServerError)
			return
		}
		var result []types.Link
		for _, item := range list {
			var timestamp int64
			if item.LastAvailability != nil {
				timestamp = item.LastAvailability.UTC().Unix()
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

		err := dbCtx.ChatLinkStore.Delete(r.Context().Value(middleware.UserIDKey).(int64), uint(linkID))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.JSONError(w, types.ErrBadRequest, "Link not found", http.StatusNotFound)
			return
		} else if err != nil {
			utils.JSONError(w, types.ErrInternalServer, "Error deleting link", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
