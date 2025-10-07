package bot

import (
	"github.com/GoBotApiOfficial/gobotapi"
	"github.com/TrackFlight/TestFlightTrackBot/internal/config"
	"github.com/TrackFlight/TestFlightTrackBot/internal/db"
	"github.com/TrackFlight/TestFlightTrackBot/internal/tor"
	"github.com/TrackFlight/TestFlightTrackBot/internal/utils"
)

type Bot struct {
	Api       *gobotapi.PollingClient
	db        *db.DB
	cfg       *config.Config
	torClient *tor.Client
	mutex     *utils.SmartMutex[int64]
}

func NewBot(cfg *config.Config, dbCtx *db.DB, torClient *tor.Client) (*Bot, error) {
	api := gobotapi.NewClient(cfg.TelegramToken)
	b := &Bot{
		Api:       api,
		db:        dbCtx,
		cfg:       cfg,
		torClient: torClient,
		mutex:     utils.NewSmartMutex[int64](),
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
