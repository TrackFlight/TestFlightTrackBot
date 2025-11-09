package testflight

import "github.com/TrackFlight/TestFlightTrackBot/internal/tor"

type Client struct {
	TorClient *tor.Client
}

func NewClient() (*Client, error) {
	return &Client{
		TorClient: &tor.Client{},
	}, nil
}
