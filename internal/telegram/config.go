package telegram

import (
	"github.com/GoBotApiOfficial/gobotapi/filters"
	"time"
)

var SupportedBotAliases = []string{
	".",
	"/",
	"!",
}

var DefaultAntiFlood = filters.AntiFlood(4, time.Second*5, time.Second*10)
