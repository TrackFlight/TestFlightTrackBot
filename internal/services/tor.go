package services

import (
	"github.com/Laky-64/gologging"
	"github.com/TrackFlight/TestFlightTrackBot/internal/tor"
	"github.com/robfig/cron/v3"
)

func startTorRotate(c *cron.Cron, t *tor.Client) error {
	_, err := c.AddFunc("*/15 * * * *", func() {
		if err := t.Refresh(); err != nil {
			gologging.ErrorF("tor refresh: %v", err)
			return
		}
	})
	return err
}
