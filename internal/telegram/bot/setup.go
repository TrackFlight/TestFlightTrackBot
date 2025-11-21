package bot

import (
	"github.com/GoBotApiOfficial/gobotapi"
	"github.com/GoBotApiOfficial/gobotapi/filters"
	"github.com/GoBotApiOfficial/gobotapi/types"
	"github.com/TrackFlight/TestFlightTrackBot/internal/db/models"
	"github.com/TrackFlight/TestFlightTrackBot/internal/telegram"
	"github.com/TrackFlight/TestFlightTrackBot/internal/telegram/handlers"
	"github.com/TrackFlight/TestFlightTrackBot/internal/telegram/handlers/admin"
)

func (b *Bot) setupHandlers() {
	// Miscellaneous commands
	b.OnCommand(
		"start",
		handlers.Start,
		filters.Private(),
	)
	b.OnCallbackQuery(
		"close",
		handlers.Close,
	)

	// Translator commands
	b.OnAdminMessage(
		admin.EditVarAction,
		filters.Private(),
		filters.Commands([]string{
			"new_var",
			"edit_var",
			"del_var",
		}, telegram.SupportedBotAliases...),
	)
	b.OnAdminCommand(
		"search_var",
		admin.SearchVar,
		filters.Private(),
	)
	b.OnAdminCommand(
		"backup",
		admin.ExecuteBackup,
		filters.Private(),
	)
	b.OnAdminCommand(
		"stats",
		admin.GetBotStats,
		filters.Private(),
	)
	b.OnAdminMessage(
		admin.RestoreBackup,
		filters.Private(),
		IsBackupFile(),
	)

	b.Api.OnMyChatMember(func(client *gobotapi.Client, update types.ChatMemberUpdated) {
		var status models.UserStatusEnum
		switch update.NewChatMember.Kind() {
		case types.TypeChatMemberBanned:
			status = models.UserStatusEnumBlockedByUser
		case types.TypeChatMemberRestricted, types.TypeChatMemberLeft:
			status = models.UserStatusEnumDeletedAccount
		default:
			status = models.UserStatusEnumReachable
		}
		_ = b.db.ChatStore.UpdateNotifiableUser(update.Chat.ID, status)
	})
}
