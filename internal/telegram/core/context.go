package core

import (
	"github.com/GoBotApiOfficial/gobotapi"
	"github.com/GoBotApiOfficial/gobotapi/methods"
	"github.com/TrackFlight/TestFlightTrackBot/internal/config"
	"github.com/TrackFlight/TestFlightTrackBot/internal/db"
	"github.com/TrackFlight/TestFlightTrackBot/internal/translator"
)

type UpdateContext struct {
	Api        *gobotapi.PollingClient
	DB         *db.DB
	Config     *config.Config
	Translator *translator.Translator
}

func NewLightContext(api *gobotapi.PollingClient, languageCode string) *UpdateContext {
	return NewUpdateContext(api, nil, nil, languageCode)
}

func NewUpdateContext(api *gobotapi.PollingClient, cfg *config.Config, dbCtx *db.DB, languageCode string) *UpdateContext {
	return &UpdateContext{
		Api:        api,
		DB:         dbCtx,
		Config:     cfg,
		Translator: translator.New(languageCode),
	}
}

func (ctx *UpdateContext) SendMessage(chatID int64, text string) error {
	return ctx.SendMessageWithKeyboard(chatID, text, nil)
}

func (ctx *UpdateContext) SendMessageWithKeyboard(chatID int64, text string, replyMarkup interface{}) error {
	_, err := ctx.Api.Invoke(
		&methods.SendMessage{
			ChatID:      chatID,
			Text:        text,
			ParseMode:   "html",
			ReplyMarkup: replyMarkup,
		},
	)
	return err
}
