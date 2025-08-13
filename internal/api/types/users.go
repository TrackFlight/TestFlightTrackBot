package types

import (
	"github.com/TrackFlight/TestFlightTrackBot/internal/db/models"
	"github.com/jackc/pgx/v5/pgtype"
)

type BaseApp struct {
	ID          int64       `json:"id"`
	Name        pgtype.Text `json:"name"`
	IconURL     pgtype.Text `json:"icon_url"`
	Description pgtype.Text `json:"description"`
	Followers   int64       `json:"followers"`
}

type App struct {
	BaseApp
	Links []Link `json:"links"`
}

type TrendingApp struct {
	BaseApp
	Links []BaseLink `json:"links"`
}

type BaseLink struct {
	ID               int64                 `json:"id"`
	URL              string                `json:"url"`
	Status           models.LinkStatusEnum `json:"status"`
	IsPublic         bool                  `json:"is_public"`
	LastAvailability int64                 `json:"last_availability"`
	LastUpdate       int64                 `json:"last_update"`
}

type Link struct {
	BaseLink
	NotifyAvailable bool `json:"notify_available"`
	NotifyClosed    bool `json:"notify_closed"`
}

type AddLinkRequest struct {
	Link            string `json:"link"`
	NotifyAvailable bool   `json:"notify_available"`
	NotifyClosed    bool   `json:"notify_closed"`
}

type EditLinkSettingsRequest struct {
	NotifyAvailable bool `json:"notify_available"`
	NotifyClosed    bool `json:"notify_closed"`
}
