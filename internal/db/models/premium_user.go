package models

type PremiumUser struct {
	ChatID int64 `gorm:"primaryKey"`
	BaseModel

	Chat Chat `gorm:"foreignKey:ChatID;references:ID;constraint:OnDelete:CASCADE"`
}
