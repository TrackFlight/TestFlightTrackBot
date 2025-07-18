package tracking

import (
	"github.com/GoBotApiOfficial/gobotapi/types"
	"github.com/Laky-64/TestFlightTrackBot/internal/telegram/core"
	"github.com/Laky-64/TestFlightTrackBot/internal/telegram/keyboards"
	"github.com/Laky-64/TestFlightTrackBot/internal/translator"
	"strconv"
)

func SendLimitReachedMessage(ctx *core.UpdateContext, message types.Message) error {
	ctx.DB.PendingStore.Remove(message.Chat.ID)
	return ctx.SendMessageWithKeyboard(
		message.Chat.ID,
		ctx.Translator.TWithData(
			translator.TrackingLimitReached,
			map[string]string{
				"MaxLinks": strconv.Itoa(ctx.Config.LimitFree),
			},
		),
		keyboards.Home(ctx.Translator),
	)
}
