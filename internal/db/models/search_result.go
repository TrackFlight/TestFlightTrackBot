package models

import "time"

type SearchResult struct {
	AppID      uint
	AppName    string
	Followers  int
	LinksCount int
	UpdatedAt  time.Time
}
