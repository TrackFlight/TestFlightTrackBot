package handlers

import (
	"github.com/GoBotApiOfficial/gobotapi/methods"
	"github.com/GoBotApiOfficial/gobotapi/types"
	"github.com/Laky-64/TestFlightTrackBot/internal/telegram/core"
)

func Close(ctx *core.UpdateContext, cb types.CallbackQuery) error {
	_, err := ctx.Api.Invoke(&methods.DeleteMessage{
		ChatID:    cb.Message.Chat.ID,
		MessageID: cb.Message.MessageID,
	})
	return err
}
