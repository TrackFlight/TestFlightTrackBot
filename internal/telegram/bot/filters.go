package bot

import (
	"github.com/GoBotApiOfficial/gobotapi"
	"github.com/GoBotApiOfficial/gobotapi/filters"
	"github.com/GoBotApiOfficial/gobotapi/types"
	"github.com/Laky-64/TestFlightTrackBot/internal/telegram"
	"github.com/Laky-64/TestFlightTrackBot/internal/telegram/core"
	"github.com/Laky-64/TestFlightTrackBot/internal/telegram/keyboards"
	"github.com/Laky-64/TestFlightTrackBot/internal/translator"
)

func buildFilter[T filters.Filterable](b *Bot, handler func(*core.UpdateContext, T) error, extraFilters ...filters.FilterOperand) func(*gobotapi.Client, T) {
	var combinedOperands []filters.FilterOperand
	combinedOperands = append(combinedOperands, filters.Private())
	combinedOperands = append(combinedOperands, extraFilters...)
	combinedOperands = append(combinedOperands, telegram.DefaultAntiFlood)
	return filters.Filter(func(client *gobotapi.Client, update T) {
		var chatID int64
		var languageHint string
		if msg, ok := any(update).(types.Message); ok {
			chatID = msg.Chat.ID
			if msg.From.ID == msg.Chat.ID {
				languageHint = msg.From.LanguageCode
			}
		} else if cb, ok := any(update).(types.CallbackQuery); ok {
			chatID = cb.Message.Chat.ID
			if cb.From.ID == cb.Message.Chat.ID {
				languageHint = cb.From.LanguageCode
			}
		}
		b.mutex.Lock(chatID)
		languageCode := b.db.ChatLinkStore.GetLanguage(chatID, languageHint)
		ctx := core.NewUpdateContext(b.Api, b.cfg, b.db, languageCode)
		err := handler(ctx, update)
		if err != nil {
			_ = ctx.SendMessageWithKeyboard(
				chatID,
				ctx.Translator.T(translator.UnknownError),
				keyboards.Home(ctx.Translator),
			)
		}
		b.mutex.Unlock(chatID)
	}, filters.And(combinedOperands...))
}

func (b *Bot) OnAdminCommand(command string, handler func(*core.UpdateContext, types.Message) error, extraFilters ...filters.FilterOperand) {
	var combinedOperands []filters.FilterOperand
	combinedOperands = append(combinedOperands, filters.UserID(b.cfg.AdminID))
	combinedOperands = append(combinedOperands, extraFilters...)
	b.OnCommand(command, handler, combinedOperands...)
}

func (b *Bot) OnCommand(command string, handler func(*core.UpdateContext, types.Message) error, extraFilters ...filters.FilterOperand) {
	var combinedOperands []filters.FilterOperand
	combinedOperands = append(combinedOperands, filters.Command(command, telegram.SupportedBotAliases...))
	combinedOperands = append(combinedOperands, extraFilters...)
	b.Api.OnMessage(
		buildFilter(
			b,
			handler,
			combinedOperands...,
		),
	)
}

func (b *Bot) OnAdminTextCommand(key translator.Key, handler func(*core.UpdateContext, types.Message) error, extraFilters ...filters.FilterOperand) {
	var combinedOperands []filters.FilterOperand
	combinedOperands = append(combinedOperands, filters.UserID(b.cfg.AdminID))
	combinedOperands = append(combinedOperands, extraFilters...)
	b.OnTextCommand(key, handler, combinedOperands...)
}

func (b *Bot) OnTextCommand(key translator.Key, handler func(*core.UpdateContext, types.Message) error, extraFilters ...filters.FilterOperand) {
	var combinedOperands []filters.FilterOperand
	combinedOperands = append(combinedOperands, matchExactTranslation(key))
	combinedOperands = append(combinedOperands, extraFilters...)
	b.Api.OnMessage(
		buildFilter(
			b,
			handler,
			combinedOperands...,
		),
	)
}

func (b *Bot) OnAdminMessage(handler func(*core.UpdateContext, types.Message) error, extraFilters ...filters.FilterOperand) {
	var combinedOperands []filters.FilterOperand
	combinedOperands = append(combinedOperands, filters.UserID(b.cfg.AdminID))
	combinedOperands = append(combinedOperands, extraFilters...)
	b.OnMessage(
		handler,
		combinedOperands...,
	)
}

func (b *Bot) OnMessage(handler func(*core.UpdateContext, types.Message) error, extraFilters ...filters.FilterOperand) {
	b.Api.OnMessage(
		buildFilter(
			b,
			handler,
			extraFilters...,
		),
	)
}

func (b *Bot) OnCallbackQuery(data string, handler func(*core.UpdateContext, types.CallbackQuery) error, extraFilters ...filters.FilterOperand) {
	var combinedOperands []filters.FilterOperand
	combinedOperands = append(combinedOperands, CallbackQueryData(data))
	combinedOperands = append(combinedOperands, extraFilters...)
	b.Api.OnCallbackQuery(
		buildFilter(
			b,
			handler,
			combinedOperands...,
		),
	)
}

func (b *Bot) IfElse(handler, handlerElse func(*core.UpdateContext, types.Message) error, filterCondition filters.FilterOperand) func(*core.UpdateContext, types.Message) error {
	return func(ctx *core.UpdateContext, message types.Message) error {
		dataFilter := &filters.DataFilter{
			Chat:    &message.Chat,
			Date:    &message.Date,
			From:    message.From,
			Client:  ctx.Api.Client,
			Message: (*types.MaybeInaccessibleMessage)(&message),
		}
		if filterCondition(dataFilter) {
			return handlerElse(ctx, message)
		}
		return handler(ctx, message)
	}
}

func (b *Bot) IsPending() filters.FilterOperand {
	return func(values *filters.DataFilter) bool {
		if b.db.PendingStore.Exists(values.Chat.ID) {
			return true
		}
		return false
	}
}

func (b *Bot) HasReachedLimit() filters.FilterOperand {
	return func(values *filters.DataFilter) bool {
		count, err := b.db.ChatLinkStore.TrackedCount(values.Chat.ID)
		if err != nil {
			return false
		}
		if count >= b.cfg.LimitFree {
			return true
		}
		return false
	}
}
