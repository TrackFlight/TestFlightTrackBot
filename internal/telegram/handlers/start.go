package handlers

import (
	"github.com/GoBotApiOfficial/gobotapi/types"
	"github.com/Laky-64/TestFlightTrackBot/internal/telegram"
	"github.com/Laky-64/TestFlightTrackBot/internal/telegram/core"
	"github.com/Laky-64/TestFlightTrackBot/internal/telegram/keyboards"
	"github.com/Laky-64/TestFlightTrackBot/internal/translator"
)

func Start(ctx *core.UpdateContext, message types.Message) error {
	return ctx.SendMessageWithKeyboard(
		message.Chat.ID,
		ctx.Translator.TWithData(
			translator.StartMessage,
			map[string]string{
				"Mention": telegram.FormatName(*message.From),
				"BotName": telegram.Mention(ctx.Api.Self()),
			},
		),
		keyboards.Home(ctx.Translator),
	)
}
