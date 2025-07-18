package models

import "time"

type TranslatorLanguage struct {
	ChatID    int64  `gorm:"primaryKey"`
	Lang      string `gorm:"primaryKey;size:5"`
	CreatedAt time.Time

	Chat Chat `gorm:"foreignKey:ChatID;references:ID;constraint:OnDelete:CASCADE"`
}
