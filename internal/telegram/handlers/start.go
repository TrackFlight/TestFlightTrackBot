package handlers

import (
	"github.com/GoBotApiOfficial/gobotapi/types"
	"github.com/TrackFlight/TestFlightTrackBot/internal/telegram"
	"github.com/TrackFlight/TestFlightTrackBot/internal/telegram/core"
	"github.com/TrackFlight/TestFlightTrackBot/internal/translator"
	"html"
)

func Start(ctx *core.UpdateContext, message types.Message) error {
	return ctx.SendMessageWithKeyboard(
		message.Chat.ID,
		ctx.Translator.TWithData(
			translator.StartMessage,
			map[string]string{
				"Mention": html.EscapeString(telegram.FormatName(*message.From)),
				"BotName": telegram.Mention(ctx.Api.Self()),
			},
		),
		&types.InlineKeyboardMarkup{
			InlineKeyboard: [][]types.InlineKeyboardButton{
				{
					{
						Text: ctx.Translator.T(translator.OpenMiniappBtn),
						WebApp: &types.WebAppInfo{
							URL: ctx.Config.MiniAppURL,
						},
					},
				},
				{
					{
						Text: ctx.Translator.T(translator.SourceCodeBtn),
						URL:  "https://github.com/TrackFlight/TestFlightTrackBot",
					},
				},
			},
		},
	)
}
