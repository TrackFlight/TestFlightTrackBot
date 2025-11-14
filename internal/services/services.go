package services

import (
	"time"

	"github.com/Laky-64/gologging"
	"github.com/TrackFlight/TestFlightTrackBot/internal/config"
	"github.com/TrackFlight/TestFlightTrackBot/internal/db"
	"github.com/TrackFlight/TestFlightTrackBot/internal/telegram/bot"
	"github.com/TrackFlight/TestFlightTrackBot/internal/testflight"
	"github.com/TrackFlight/TestFlightTrackBot/internal/utils"
	"github.com/robfig/cron/v3"
)

func StartAll(
	c *cron.Cron,
	b *bot.Bot,
	cfg *config.Config,
	dbCtx *db.DB,
	tfClient *testflight.Client,
) {
	rateLimit := utils.NewRateLimiter(
		MaxMessagesPerSecond,
		time.Second,
	)

	if err := startTestflight(c, rateLimit, b, cfg, dbCtx, tfClient); err != nil {
		gologging.FatalF("start testflight service: %v", err)
	}
	if err := startTorRotate(c, tfClient.TorClient); err != nil {
		gologging.FatalF("start tor rotate service: %v", err)
	}
	if err := startDbBackup(c, b, cfg); err != nil {
		gologging.FatalF("start db backup service: %v", err)
	}
	if err := startWeeklyHighLights(c, rateLimit, b, cfg, dbCtx, tfClient.TorClient); err != nil {
		gologging.FatalF("start weekly highlights service: %v", err)
	}
}
