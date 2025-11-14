package core

import (
	"github.com/GoBotApiOfficial/gobotapi"
	"github.com/GoBotApiOfficial/gobotapi/methods"
	rawTypes "github.com/GoBotApiOfficial/gobotapi/types/raw"
	"github.com/TrackFlight/TestFlightTrackBot/internal/config"
	"github.com/TrackFlight/TestFlightTrackBot/internal/db"
	"github.com/TrackFlight/TestFlightTrackBot/internal/tor"
	"github.com/TrackFlight/TestFlightTrackBot/internal/translator"
)

type UpdateContext struct {
	Api        *gobotapi.PollingClient
	DB         *db.DB
	TorClient  *tor.Client
	Config     *config.Config
	Translator *translator.Translator
}

func NewLightContext(api *gobotapi.PollingClient, languageCode string) *UpdateContext {
	return NewUpdateContext(api, nil, nil, nil, languageCode)
}

func NewUpdateContext(api *gobotapi.PollingClient, cfg *config.Config, torClient *tor.Client, dbCtx *db.DB, languageCode string) *UpdateContext {
	return &UpdateContext{
		Api:        api,
		DB:         dbCtx,
		Config:     cfg,
		Translator: translator.New(languageCode),
		TorClient:  torClient,
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

func (ctx *UpdateContext) SendPhotoWithKeyboard(chatID int64, text string, photo rawTypes.InputFile, replyMarkup interface{}) (*rawTypes.Result, error) {
	return ctx.Api.Invoke(
		&methods.SendPhoto{
			ChatID:      chatID,
			Caption:     text,
			Photo:       photo,
			ParseMode:   "html",
			ReplyMarkup: replyMarkup,
		},
	)
}
