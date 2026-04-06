package handlers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/GoBotApiOfficial/gobotapi/methods"
	"github.com/GoBotApiOfficial/gobotapi/types"
	"github.com/TrackFlight/TestFlightTrackBot/internal/telegram/core"
	"github.com/TrackFlight/TestFlightTrackBot/internal/translator"
)

func Mute(ctx *core.UpdateContext, cb types.CallbackQuery) error {
	parts := strings.SplitN(cb.Data, ":", 4)
	if len(parts) != 4 {
		return nil
	}
	notifType, linkIdStr, linkCode := parts[1], parts[2], parts[3]
	_, _ = ctx.Api.Invoke(&methods.AnswerCallbackQuery{CallbackQueryID: cb.ID})
	_, _ = ctx.Api.Invoke(&methods.EditMessageReplyMarkup{
		ChatID:    cb.Message.Chat.ID,
		MessageID: cb.Message.MessageID,
		ReplyMarkup: &types.InlineKeyboardMarkup{
			InlineKeyboard: [][]types.InlineKeyboardButton{
				{
					{
						Text:         ctx.Translator.T(translator.MuteNotificationsConfirmBtn),
						CallbackData: fmt.Sprintf("mute_y:%s", linkIdStr),
					},
					{
						Text:         ctx.Translator.T(translator.MuteNotificationsCancelBtn),
						CallbackData: fmt.Sprintf("mute_n:%s:%s:%s", notifType, linkIdStr, linkCode),
					},
				},
			},
		},
	})
	return nil
}

func MuteConfirm(ctx *core.UpdateContext, cb types.CallbackQuery) error {
	parts := strings.SplitN(cb.Data, ":", 2)
	if len(parts) != 2 {
		return nil
	}
	linkId, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return nil
	}
	_ = ctx.DB.ChatStore.UpdateNotificationSettings(
		cb.Message.Chat.ID,
		linkId,
		false,
		false,
		ctx.Config.LimitFree,
	)
	_, _ = ctx.Api.Invoke(&methods.AnswerCallbackQuery{
		CallbackQueryID: cb.ID,
		Text:            ctx.Translator.T(translator.MuteNotificationsSuccess),
	})
	_, _ = ctx.Api.Invoke(&methods.EditMessageReplyMarkup{
		ChatID:    cb.Message.Chat.ID,
		MessageID: cb.Message.MessageID,
	})
	return nil
}

func MuteCancel(ctx *core.UpdateContext, cb types.CallbackQuery) error {
	parts := strings.SplitN(cb.Data, ":", 4)
	if len(parts) != 4 {
		return nil
	}
	notifType, linkIdStr, linkCode := parts[1], parts[2], parts[3]
	muteData := fmt.Sprintf("mute:%s:%s:%s", notifType, linkIdStr, linkCode)
	muteBtn := types.InlineKeyboardButton{
		Text:         ctx.Translator.T(translator.MuteNotificationsBtn),
		CallbackData: muteData,
	}
	var rows [][]types.InlineKeyboardButton
	if notifType == "o" {
		linkURL := "https://testflight.apple.com/join/" + linkCode
		rows = [][]types.InlineKeyboardButton{
			{
				{
					Text: ctx.Translator.T(translator.JoinBetaBtn),
					URL:  linkURL,
				},
				muteBtn,
			},
		}
	} else {
		rows = [][]types.InlineKeyboardButton{
			{muteBtn},
		}
	}
	_, _ = ctx.Api.Invoke(&methods.AnswerCallbackQuery{CallbackQueryID: cb.ID})
	_, _ = ctx.Api.Invoke(&methods.EditMessageReplyMarkup{
		ChatID:    cb.Message.Chat.ID,
		MessageID: cb.Message.MessageID,
		ReplyMarkup: &types.InlineKeyboardMarkup{
			InlineKeyboard: rows,
		},
	})
	return nil
}
