package tracking

import (
	"github.com/GoBotApiOfficial/gobotapi/types"
	"github.com/Laky-64/TestFlightTrackBot/internal/telegram/core"
	"github.com/Laky-64/TestFlightTrackBot/internal/translator"
)

func Start(ctx *core.UpdateContext, message types.Message) error {
	if err := ctx.DB.PendingStore.Add(message.From.ID); err != nil {
		return err
	}
	return ctx.SendMessageWithKeyboard(
		message.Chat.ID,
		ctx.Translator.T(translator.TrackingWaiting),
		&types.ReplyKeyboardRemove{
			RemoveKeyboard: true,
		},
	)
}
