package main

import (
	"context"
	"github.com/Laky-64/TestFlightTrackBot/internal/api"
	"github.com/Laky-64/TestFlightTrackBot/internal/config"
	"github.com/Laky-64/TestFlightTrackBot/internal/db"
	"github.com/Laky-64/TestFlightTrackBot/internal/services"
	"github.com/Laky-64/TestFlightTrackBot/internal/telegram/bot"
	"github.com/Laky-64/TestFlightTrackBot/internal/testflight"
	"github.com/Laky-64/gologging"
	"os"
	"os/signal"
	"syscall"
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

	b, err := bot.NewBot(cfg, dbCtx)
	if err != nil {
		gologging.Fatal(err)
	}

	tfClient, err := testflight.NewClient()
	if err != nil {
		gologging.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	services.StartAll(ctx, b, cfg, dbCtx, tfClient)

	api.Start(dbCtx, cfg)

	if err = b.Start(); err != nil {
		gologging.Fatal(err)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
}
