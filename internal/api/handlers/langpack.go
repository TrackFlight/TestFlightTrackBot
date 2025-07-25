package handlers

import (
	"encoding/json"
	"github.com/Laky-64/TestFlightTrackBot/internal/api/middleware"
	"github.com/Laky-64/TestFlightTrackBot/internal/api/types"
	"github.com/Laky-64/TestFlightTrackBot/internal/api/utils"
	"github.com/Laky-64/TestFlightTrackBot/internal/db"
	"github.com/Laky-64/TestFlightTrackBot/internal/translator"
	"golang.org/x/text/language"
	"net/http"
)

func GetLangPack(dbCtx *db.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		langCode := r.URL.Query().Get("lang_code")
		langCodeHint := r.URL.Query().Get("lang_code_hint")
		if _, err := language.Parse(langCode); err != nil && langCode != "" {
			utils.JSONError(w, types.ErrInvalidLanguageCode, "Invalid language code", http.StatusBadRequest)
			return
		}
		if _, err := language.Parse(langCodeHint); err != nil && langCodeHint != "" || langCode == "" && langCodeHint == "" {
			utils.JSONError(w, types.ErrInvalidLanguageCode, "Invalid language hint", http.StatusBadRequest)
			return
		}

		if len(langCode) == 0 {
			langCode, _ = dbCtx.ChatStore.GetLanguage(
				r.Context().Value(middleware.UserIDKey).(int64),
				r.URL.Query().Get("languageHint"),
			)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		marshal, _ := json.Marshal(translator.LangPack(langCode))
		_, _ = w.Write(marshal)
	}
}
