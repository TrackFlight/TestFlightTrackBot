package types

type App struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	IconURL     string `json:"icon_url"`
	Description string `json:"description"`
	Links       []Link `json:"links"`
}

type Link struct {
	ID               int64  `json:"id"`
	DisplayName      string `json:"display_name"`
	Status           string `json:"status"`
	LastAvailability int64  `json:"last_availability"`
	LastUpdate       int64  `json:"last_update"`
}

type AddLinkRequest struct {
	ID   int64  `json:"id"`
	Link string `json:"link"`
}
