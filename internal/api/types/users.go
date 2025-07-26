package types

import (
	"github.com/Laky-64/TestFlightTrackBot/internal/db/models"
	"github.com/jackc/pgx/v5/pgtype"
)

type App struct {
	ID          pgtype.Int8 `json:"id"`
	Name        pgtype.Text `json:"name"`
	IconURL     pgtype.Text `json:"icon_url"`
	Description pgtype.Text `json:"description"`
	Links       []Link      `json:"links"`
}

type Link struct {
	ID               int64                 `json:"id"`
	DisplayName      string                `json:"display_name"`
	Status           models.LinkStatusEnum `json:"status"`
	LastAvailability int64                 `json:"last_availability"`
	LastUpdate       int64                 `json:"last_update"`
}

type AddLinkRequest struct {
	ID   int64  `json:"id"`
	Link string `json:"link"`
}
