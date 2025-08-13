package services

import (
	"context"
	"github.com/TrackFlight/TestFlightTrackBot/internal/config"
	"github.com/TrackFlight/TestFlightTrackBot/internal/db"
	"github.com/TrackFlight/TestFlightTrackBot/internal/telegram/bot"
	"github.com/TrackFlight/TestFlightTrackBot/internal/testflight"
	"github.com/TrackFlight/TestFlightTrackBot/internal/utils"
	"time"
)

func StartAll(
	ctx context.Context,
	b *bot.Bot,
	cfg *config.Config,
	dbCtx *db.DB,
	tfClient *testflight.Client,
) {
	rateLimit := utils.NewRateLimiter(
		MaxMessagesPerSecond,
		time.Second,
	)

	go startTestflight(ctx, rateLimit, b, cfg, dbCtx, tfClient)
	go startTorRotate(ctx, tfClient.TorClient)
	go startDbBackup(ctx, b, cfg)
}
