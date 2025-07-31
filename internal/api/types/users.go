package types

import (
	"github.com/Laky-64/TestFlightTrackBot/internal/db/models"
	"github.com/jackc/pgx/v5/pgtype"
)

type App struct {
	ID          int64       `json:"id"`
	Name        pgtype.Text `json:"name"`
	IconURL     pgtype.Text `json:"icon_url"`
	Description pgtype.Text `json:"description"`
	Followers   int64       `json:"followers"`
	Links       []Link      `json:"links"`
}

type Link struct {
	ID               int64                 `json:"id"`
	URL              string                `json:"url"`
	Status           models.LinkStatusEnum `json:"status"`
	LastAvailability int64                 `json:"last_availability"`
	LastUpdate       int64                 `json:"last_update"`
}

type AddLinkRequest struct {
	ID   int64  `json:"id"`
	Link string `json:"link"`
}
