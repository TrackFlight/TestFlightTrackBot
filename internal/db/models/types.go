package models

import "database/sql/driver"

type LinkStatus string

const (
	StatusAvailable LinkStatus = "available"
	StatusFull      LinkStatus = "full"
	StatusClosed    LinkStatus = "closed"
	StatusInvalid   LinkStatus = "invalid"
)

func (ns LinkStatus) Value() (driver.Value, error) {
	if ns == "" {
		return nil, nil
	}
	return string(ns), nil
}
