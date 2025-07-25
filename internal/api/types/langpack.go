package types

import "github.com/Laky-64/TestFlightTrackBot/internal/translator"

type LangPack struct {
	LangCode string                    `json:"lang_code"`
	Strings  map[translator.Key]string `json:"strings"`
}
