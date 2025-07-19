package models

type BaseNotificationChat struct {
	ChatID  int64
	Lang    string
	LinkURL string
}

type NotificationChat struct {
	BaseNotificationChat
	AppName string
	Status  LinkStatus
}

type NotificationRequest struct {
	LinkID uint
	Status LinkStatus
}
