package models

type ChatLink struct {
	ChatID             int64      `gorm:"primaryKey"`
	LinkID             uint       `gorm:"primaryKey"`
	AllowOpened        bool       `gorm:"default:true;not null;index"`
	AllowClosed        bool       `gorm:"default:false;not null;index"`
	LastNotifiedStatus LinkStatus `gorm:"type:link_status_enum;index"`
	BaseModel

	Chat Chat `gorm:"foreignKey:ChatID;references:ID;constraint:OnDelete:CASCADE"`
	Link Link `gorm:"foreignKey:LinkID;references:ID;constraint:OnDelete:CASCADE"`
}
