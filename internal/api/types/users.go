package types

type Link struct {
	ID               int64  `json:"id"`
	AppID            int64  `json:"app_id"`
	Tag              string `json:"tag"`
	AppName          string `json:"app_name"`
	IconURL          string `json:"icon_url"`
	Description      string `json:"description"`
	Status           string `json:"status"`
	LastAvailability int64  `json:"last_availability"`
}

type NewLink struct {
	ID   int64  `json:"id"`
	Link string `json:"link"`
}
