package types

import "github.com/TrackFlight/TestFlightTrackBot/internal/translator"

type LangPack struct {
	LangCode string                    `json:"lang_code"`
	Strings  map[translator.Key]string `json:"strings"`
}
