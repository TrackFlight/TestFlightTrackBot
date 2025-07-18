package tracking

import (
	"github.com/GoBotApiOfficial/gobotapi/types"
	"github.com/Laky-64/TestFlightTrackBot/internal/telegram/core"
	"github.com/Laky-64/TestFlightTrackBot/internal/telegram/keyboards"
	"github.com/Laky-64/TestFlightTrackBot/internal/translator"
)

func Cancel(ctx *core.UpdateContext, message types.Message) error {
	if ctx.DB.PendingStore.Remove(message.Chat.ID) {
		return ctx.SendMessageWithKeyboard(
			message.Chat.ID,
			ctx.Translator.T(translator.TrackingCancelled),
			keyboards.Home(ctx.Translator),
		)
	}
	return nil
}
