package bot

import (
	"github.com/GoBotApiOfficial/gobotapi/filters"
	"github.com/GoBotApiOfficial/gobotapi/types"
	"github.com/Laky-64/TestFlightTrackBot/internal/telegram"
	"github.com/Laky-64/TestFlightTrackBot/internal/testflight"
	"github.com/Laky-64/TestFlightTrackBot/internal/translator"
	"strings"
)

func IsCommand() filters.FilterOperand {
	return filters.Command("[A-Za-z0-9]+", telegram.SupportedBotAliases...)
}

func IsTextCommand() filters.FilterOperand {
	var stringPatterns []string
	for _, key := range translator.TKeys() {
		if strings.HasSuffix(string(key), "_btn") {
			stringPatterns = append(stringPatterns, translator.TAll(key)...)
		}
	}
	return func(df *filters.DataFilter) bool {
		message, ok := df.RawUpdate.(types.Message)
		if !ok {
			return false
		}
		for _, pattern := range stringPatterns {
			if message.Text == pattern {
				return true
			}
		}
		return false
	}
}

func IsTestFlightLink() filters.FilterOperand {
	return func(df *filters.DataFilter) bool {
		message, ok := df.RawUpdate.(types.Message)
		if !ok {
			return false
		}
		return testflight.RegexLink.MatchString(message.Text)
	}
}

func CallbackQueryData(data string) filters.FilterOperand {
	return func(df *filters.DataFilter) bool {
		query, ok := df.RawUpdate.(types.CallbackQuery)
		if !ok {
			return false
		}
		return strings.HasPrefix(query.Data, data)
	}
}

func matchExactTranslation(key translator.Key) filters.FilterOperand {
	texts := translator.TAll(key)
	return func(df *filters.DataFilter) bool {
		message, ok := df.RawUpdate.(types.Message)
		if !ok {
			return false
		}
		for _, txt := range texts {
			if message.Text == txt {
				return true
			}
		}
		return false
	}
}
