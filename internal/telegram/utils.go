package telegram

import (
	"fmt"
	"html"
	"strings"

	"github.com/GoBotApiOfficial/gobotapi/types"
	"github.com/TrackFlight/TestFlightTrackBot/internal/db/models"
)

func FormatName(user types.User) string {
	if user.LastName == "" {
		return user.FirstName
	}
	return user.FirstName + " " + user.LastName
}

func Mention(user types.User) string {
	return fmt.Sprintf(
		`<a href="tg://user?id=%d">%s</a>`,
		user.ID,
		html.EscapeString(FormatName(user)),
	)
}

func MapErrorToUserStatus(err error) models.UserStatusEnum {
	if strings.Contains(err.Error(), "bot was blocked by the user") {
		return models.UserStatusEnumBlockedByUser
	}
	if strings.Contains(err.Error(), "user is deactivated") {
		return models.UserStatusEnumDeletedAccount
	}
	return models.UserStatusEnumReachable
}
