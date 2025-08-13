package bot

import (
	"github.com/GoBotApiOfficial/gobotapi/filters"
	"github.com/Laky-64/TestFlightTrackBot/internal/telegram"
	"github.com/Laky-64/TestFlightTrackBot/internal/telegram/handlers"
	"github.com/Laky-64/TestFlightTrackBot/internal/telegram/handlers/admin"
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
	b.OnAdminMessage(
		admin.RestoreBackup,
		filters.Private(),
		IsBackupFile(),
	)
}
