package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Laky-64/gologging"
	"github.com/TrackFlight/TestFlightTrackBot/internal/api"
	"github.com/TrackFlight/TestFlightTrackBot/internal/config"
	"github.com/TrackFlight/TestFlightTrackBot/internal/db"
	"github.com/TrackFlight/TestFlightTrackBot/internal/services"
	"github.com/TrackFlight/TestFlightTrackBot/internal/telegram/bot"
	"github.com/TrackFlight/TestFlightTrackBot/internal/testflight"
	"github.com/robfig/cron/v3"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		gologging.Fatal(err)
	}

	dbCtx, err := db.NewDB(cfg)
	if err != nil {
		gologging.Fatal(err)
	}

	tfClient, err := testflight.NewClient()
	if err != nil {
		gologging.Fatal(err)
	}

	b, err := bot.NewBot(cfg, dbCtx, tfClient.TorClient)
	if err != nil {
		gologging.Fatal(err)
	}

	c := cron.New()
	services.StartAll(c, b, cfg, dbCtx, tfClient)
	defer c.Stop()

	api.Start(dbCtx, cfg)

	if err = b.Start(); err != nil {
		gologging.Fatal(err)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
}
