package utils

import (
	"encoding/json"
	"fmt"
	"github.com/TrackFlight/TestFlightTrackBot/internal/api/types"
	"net/http"
)

func JSONErrorFloodWait(w http.ResponseWriter, seconds int) {
	rawJSONError(w, types.ErrFloodWait, fmt.Sprintf("A wait of %d seconds is required", seconds), http.StatusTooManyRequests, seconds)
}

func JSONError(w http.ResponseWriter, code, message string, status int) {
	rawJSONError(w, code, message, status, 0)
}

func rawJSONError(w http.ResponseWriter, code, message string, status, seconds int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(types.ErrorResponse{
		Code:    code,
		Message: message,
		Seconds: seconds,
	})
}
