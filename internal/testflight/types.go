package testflight

import "github.com/TrackFlight/TestFlightTrackBot/internal/db/models"

type Result struct {
	Error       error
	AppName     string
	IconURL     string
	Description string
	IsPublic    bool
	Status      models.LinkStatusEnum
}
