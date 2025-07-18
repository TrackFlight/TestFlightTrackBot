package services

import (
	"context"
	"github.com/Laky-64/TestFlightTrackBot/internal/tor"
	"github.com/Laky-64/gologging"
	"time"
)

func startTorRotate(ctx context.Context, t *tor.Client) {
	ticker := time.NewTicker(15 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := t.Refresh(); err != nil {
				gologging.ErrorF("tor refresh: %v", err)
				return
			}
		}
	}
}
