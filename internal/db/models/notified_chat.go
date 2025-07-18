package models

type NotifiedChat struct {
	ChatID  int64
	Lang    string
	LinkURL string
	AppName string
	Status  LinkStatus
}
