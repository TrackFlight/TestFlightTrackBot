package admin

import (
	"fmt"
	"github.com/GoBotApiOfficial/gobotapi/methods"
	"github.com/GoBotApiOfficial/gobotapi/types"
	"github.com/TrackFlight/TestFlightTrackBot/internal/db"
	"github.com/TrackFlight/TestFlightTrackBot/internal/telegram/core"
	"time"
)

func ExecuteBackup(ctx *core.UpdateContext, message types.Message) error {
	backup, err := db.ExecuteBackup(ctx.Config)
	if err != nil {
		return err
	}
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("backup_%s.sql", timestamp)
	_, err = ctx.Api.Invoke(&methods.SendDocument{
		ChatID: message.Chat.ID,
		Document: types.InputBytes{
			Name: filename,
			Data: backup,
		},
	})
	return err
}

func RestoreBackup(ctx *core.UpdateContext, message types.Message) error {
	bytes, err := ctx.Api.DownloadBytes(message.Document.FileID, nil)
	if err != nil {
		return err
	}
	err = db.RestoreBackup(ctx.Config, bytes)
	if err != nil {
		return err
	}
	return ctx.SendMessage(
		message.Chat.ID,
		"<b>✅ Backup restored successfully!</b>",
	)
}
