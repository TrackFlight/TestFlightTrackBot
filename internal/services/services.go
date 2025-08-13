package services

import (
	"context"
	"github.com/Laky-64/TestFlightTrackBot/internal/config"
	"github.com/Laky-64/TestFlightTrackBot/internal/db"
	"github.com/Laky-64/TestFlightTrackBot/internal/telegram/bot"
	"github.com/Laky-64/TestFlightTrackBot/internal/testflight"
	"github.com/Laky-64/TestFlightTrackBot/internal/utils"
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
