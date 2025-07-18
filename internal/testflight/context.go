package testflight

import "github.com/Laky-64/TestFlightTrackBot/internal/tor"

type Client struct {
	userAgents []string
	TorClient  *tor.Client
}

func NewClient() (*Client, error) {
	agents, err := loadUserAgents()
	if err != nil {
		return nil, err
	}

	return &Client{
		userAgents: agents,
		TorClient:  &tor.Client{},
	}, nil
}
