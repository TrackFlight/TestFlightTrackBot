package tracking

import (
	"github.com/GoBotApiOfficial/gobotapi/types"
	"github.com/Laky-64/TestFlightTrackBot/internal/telegram/core"
)

func TrackLink(ctx *core.UpdateContext, message types.Message) error {
	ctx.DB.PendingStore.Remove(message.Chat.ID)
	if linkID, err := ctx.DB.LinkStore.Add(message.Text); err != nil {
		return err
	} else {
		return addTrackingLink(ctx, message.Chat.ID, linkID)
	}
}

func SearchLink(ctx *core.UpdateContext, message types.Message) error {
	//results, err := ctx.DB.AppStore.Search(message.Text)
	//if err != nil {
	//	return err
	//}
	//if len(results) == 0 {
	//	return ctx.SendMessage(
	//		message.Chat.ID,
	//		ctx.Translator.T(translator.NoAppsFound),
	//	)
	//}
	//var resultsText strings.Builder
	//for _, app := range results {
	//	resultsText.WriteString(
	//		ctx.Translator.TWithData(
	//			translator.FoundAppsTemplate,
	//			map[string]interface{}{
	//				"AppName":    app.AppName,
	//				"LinksCount": app.LinksCount,
	//				"Followers":  app.Followers,
	//				"LastUpdate": utils.ReadableDateDifference(
	//					ctx.Translator,
	//					app.UpdatedAt,
	//					time.Now(),
	//				),
	//			},
	//		),
	//	)
	//	resultsText.WriteString("\n\n")
	//}
	//return ctx.SendMessage(
	//	message.Chat.ID,
	//	resultsText.String(),
	//)
	return nil
}
