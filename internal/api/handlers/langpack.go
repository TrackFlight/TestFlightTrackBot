package handlers

import (
	"encoding/json"
	"github.com/Laky-64/TestFlightTrackBot/internal/api/middleware"
	"github.com/Laky-64/TestFlightTrackBot/internal/db"
	"github.com/Laky-64/TestFlightTrackBot/internal/translator"
	"net/http"
)

func GetStrings(dbCtx *db.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		language := r.URL.Query().Get("language")

		if len(language) == 0 {
			language, _ = dbCtx.ChatStore.GetLanguage(
				r.Context().Value(middleware.UserIDKey).(int64),
				r.URL.Query().Get("languageHint"),
			)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		marshal, _ := json.Marshal(translator.LangPack(language))
		_, _ = w.Write(marshal)
	}
}
