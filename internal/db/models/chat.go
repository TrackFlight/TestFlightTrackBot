package models

type Chat struct {
	ID   int64  `gorm:"primaryKey"`
	Lang string `gorm:"size:5;not null;default:'en'"`
	BaseModel
}
