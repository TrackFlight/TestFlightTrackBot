package types

type Link struct {
	ID               uint   `json:"id"`
	Tag              string `json:"tag"`
	AppName          string `json:"app_name"`
	IconURL          string `json:"icon_url"`
	Description      string `json:"description"`
	Status           string `json:"status"`
	LastAvailability int64  `json:"last_availability"`
}
