package telegram

import (
	"fmt"

	"github.com/GoBotApiOfficial/gobotapi/types"
	"github.com/TrackFlight/TestFlightTrackBot/internal/translator"
)

// Premium Emojis
const (
	MuteIconEmojiID   = "5244807637157029775"
	ChainsIconEmojiID = "5375129357373165375"
	CrossIconEmojiID  = "5465665476971471368"
)

func NotificationKeyboard(t *translator.Translator, linkID int64, link string, opened bool) *types.InlineKeyboardMarkup {
	notifType := "c"
	if opened {
		notifType = "o"
	}
	keyboard := &types.InlineKeyboardMarkup{
		InlineKeyboard: [][]types.InlineKeyboardButton{{}},
	}
	if opened {
		keyboard.InlineKeyboard[0] = append(keyboard.InlineKeyboard[0], types.InlineKeyboardButton{
			Text:              t.T(translator.JoinBetaBtn),
			IconCustomEmojiID: ChainsIconEmojiID,
			URL:               link,
		})
	}
	keyboard.InlineKeyboard[0] = append(keyboard.InlineKeyboard[0], types.InlineKeyboardButton{
		Text:              t.T(translator.MuteNotificationsBtn),
		IconCustomEmojiID: MuteIconEmojiID,
		CallbackData:      fmt.Sprintf("mute:%s:%d:%s", notifType, linkID, link),
	})
	return keyboard
}
