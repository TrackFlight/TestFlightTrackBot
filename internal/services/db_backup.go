package services

import (
	"fmt"
	"time"

	"github.com/GoBotApiOfficial/gobotapi/methods"
	"github.com/GoBotApiOfficial/gobotapi/types"
	"github.com/Laky-64/gologging"
	"github.com/TrackFlight/TestFlightTrackBot/internal/config"
	"github.com/TrackFlight/TestFlightTrackBot/internal/db"
	"github.com/TrackFlight/TestFlightTrackBot/internal/telegram/bot"
	"github.com/robfig/cron/v3"
)

func startDbBackup(c *cron.Cron, b *bot.Bot, cfg *config.Config) error {
	_, err := c.AddFunc("0 */6 * * *", func() {
		backup, err := db.ExecuteBackup(cfg)
		if err != nil {
			gologging.Error("Error executing backup", err)
			return
		}
		timestamp := time.Now().Format("2006-01-02_15-04-05")
		filename := fmt.Sprintf("backup_%s.sql", timestamp)
		_, err = b.Api.Invoke(
			&methods.SendDocument{
				ChatID: cfg.BackupChatID,
				Document: types.InputBytes{
					Name: filename,
					Data: backup,
				},
			},
		)
	})
	return err
}
