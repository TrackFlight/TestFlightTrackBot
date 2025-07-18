package telegram

import (
	"fmt"
	"github.com/GoBotApiOfficial/gobotapi/types"
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
		FormatName(user),
	)
}
