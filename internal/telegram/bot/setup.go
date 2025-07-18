package bot

import (
	"github.com/GoBotApiOfficial/gobotapi/filters"
	"github.com/Laky-64/TestFlightTrackBot/internal/telegram"
	"github.com/Laky-64/TestFlightTrackBot/internal/telegram/handlers"
	"github.com/Laky-64/TestFlightTrackBot/internal/telegram/handlers/admin"
	"github.com/Laky-64/TestFlightTrackBot/internal/telegram/handlers/tracking"
	"github.com/Laky-64/TestFlightTrackBot/internal/translator"
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

	// Tracking commands
	b.OnTextCommand(
		translator.StartTrackingBtn,
		b.IfElse(
			tracking.Start,
			tracking.SendLimitReachedMessage,
			b.HasReachedLimit(),
		),
	)
	b.OnCommand(
		"cancel",
		tracking.Cancel,
	)
	b.OnMessage(
		b.IfElse(
			tracking.TrackLink,
			tracking.SendLimitReachedMessage,
			b.HasReachedLimit(),
		),
		IsTestFlightLink(),
	)
	b.OnMessage(
		tracking.SearchLink,
		filters.Not(
			filters.Or(
				IsCommand(),
				IsTestFlightLink(),
				IsTextCommand(),
			),
		),
		b.IsPending(),
	)
}
