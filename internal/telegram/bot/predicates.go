package bot

import (
	"github.com/GoBotApiOfficial/gobotapi/filters"
	"github.com/GoBotApiOfficial/gobotapi/types"
	"github.com/Laky-64/TestFlightTrackBot/internal/translator"
	"strings"
)

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
