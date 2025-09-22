package telegram

import (
	"time"

	"github.com/GoBotApiOfficial/gobotapi/filters"
)

var SupportedBotAliases = []string{
	".",
	"/",
	"!",
}

var DefaultAntiFlood = filters.AntiFlood(4, time.Second*5, time.Second*10)
