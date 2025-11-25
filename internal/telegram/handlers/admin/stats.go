package admin

import (
	"fmt"
	"strconv"

	"github.com/GoBotApiOfficial/gobotapi/types"
	"github.com/TrackFlight/TestFlightTrackBot/internal/telegram/core"
	"github.com/TrackFlight/TestFlightTrackBot/internal/translator"
)

func GetBotStats(ctx *core.UpdateContext, message types.Message) error {
	stats, err := ctx.DB.AdminStore.GetBotStats()
	if err != nil || stats == nil {
		return err
	}
	statsMessage := "<b>🤖 Bot Statistics:</b>\n" +
		" ├ 👤 <b>Users:</b> <code>" + fmt.Sprintf("%d/%d", stats.ActiveChats, stats.TotalChats) + "</code>\n" +
		" ├ ⛔️ <b>Blocked Chats:</b> <code>" + strconv.Itoa(int(stats.BlockedChats)) + "</code>\n" +
		" ├ 🔗 <b>Links:</b> <code>" + fmt.Sprintf("%d/%d", stats.TrackedLinks, stats.TotalLinks) + "</code>\n" +
		" ╰ 🧅 <b>TOR Instances:</b> <code>" + strconv.Itoa(ctx.TorClient.InstanceCount()) + "</code>\n"
	return ctx.SendMessageWithKeyboard(
		message.Chat.ID,
		statsMessage,
		&types.InlineKeyboardMarkup{
			InlineKeyboard: [][]types.InlineKeyboardButton{
				{
					{
						Text:         ctx.Translator.T(translator.CloseBtn),
						CallbackData: "close",
					},
				},
			},
		},
	)
}
