package models

import "time"

type PendingTrack struct {
	ChatID    int64     `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"index"`

	Chat Chat `gorm:"foreignKey:ChatID;references:ID;constraint:OnDelete:CASCADE"`
}
