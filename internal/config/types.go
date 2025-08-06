package config

type Config struct {
	AdminID int64

	DBUser     string
	DBPassword string
	DBName     string

	TelegramToken string

	LimitFree         int64
	LimitPremium      int64
	MaxFollowingLinks int64

	PublicLinkMinUsers int64
}
