package bot

import (
	"github.com/GoBotApiOfficial/gobotapi"
	"github.com/Laky-64/TestFlightTrackBot/internal/config"
	"github.com/Laky-64/TestFlightTrackBot/internal/db"
	"github.com/Laky-64/TestFlightTrackBot/internal/utils"
)

type Bot struct {
	Api   *gobotapi.PollingClient
	db    *db.DB
	cfg   *config.Config
	mutex *utils.SmartMutex[int64]
}

func NewBot(cfg *config.Config, dbCtx *db.DB) (*Bot, error) {
	api := gobotapi.NewClient(cfg.TelegramToken)
	b := &Bot{
		Api:   api,
		db:    dbCtx,
		cfg:   cfg,
		mutex: utils.NewSmartMutex[int64](),
	}
	b.setupHandlers()
	return b, nil
}

func (b *Bot) Start() error {
	if err := b.Api.Start(); err != nil {
		return err
	}
	return nil
}
