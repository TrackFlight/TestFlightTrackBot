package tracking

import (
	"github.com/Laky-64/TestFlightTrackBot/internal/telegram/core"
	"github.com/Laky-64/TestFlightTrackBot/internal/telegram/keyboards"
	"github.com/Laky-64/TestFlightTrackBot/internal/translator"
)

func addTrackingLink(ctx *core.UpdateContext, chatID int64, linkID uint) error {
	if ok, err := ctx.DB.ChatLinkStore.Track(chatID, linkID); err != nil {
		return err
	} else if !ok {
		return ctx.SendMessageWithKeyboard(
			chatID,
			ctx.Translator.T(translator.TrackingAlreadyAdded),
			keyboards.Home(ctx.Translator),
		)
	}
	return ctx.SendMessageWithKeyboard(
		chatID,
		ctx.Translator.T(translator.TrackingStarted),
		keyboards.Home(ctx.Translator),
	)
}
