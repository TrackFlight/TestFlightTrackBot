package models

import "time"

type TrackingLink struct {
	ID               uint
	AppName          string
	IconURL          string
	Description      string
	Status           LinkStatus
	LastAvailability *time.Time
}
