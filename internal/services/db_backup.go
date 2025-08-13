package services

import (
	"context"
	"fmt"
	"github.com/GoBotApiOfficial/gobotapi/methods"
	"github.com/GoBotApiOfficial/gobotapi/types"
	"github.com/Laky-64/gologging"
	"github.com/TrackFlight/TestFlightTrackBot/internal/config"
	"github.com/TrackFlight/TestFlightTrackBot/internal/db"
	"github.com/TrackFlight/TestFlightTrackBot/internal/telegram/bot"
	"time"
)

func startDbBackup(ctx context.Context, b *bot.Bot, cfg *config.Config) {
	ticker := time.NewTicker(6 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			backup, err := db.ExecuteBackup(cfg)
			if err != nil {
				gologging.Error("Error executing backup", err)
				continue
			}
			timestamp := time.Now().Format("2006-01-02_15-04-05")
			filename := fmt.Sprintf("backup_%s.sql", timestamp)
			_, err = b.Api.Invoke(&methods.SendDocument{
				ChatID: cfg.BackupChatID,
				Document: types.InputBytes{
					Name: filename,
					Data: backup,
				},
			})
		}
	}
}
