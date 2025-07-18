package models

type App struct {
	ID          uint   `gorm:"primaryKey;not null"`
	AppName     string `gorm:"type:varchar(255);not null;unique;index"`
	IconURL     string `gorm:"type:varchar(255);not null;"`
	Description string `gorm:"type:text;not null;"`
	BaseModel
}
