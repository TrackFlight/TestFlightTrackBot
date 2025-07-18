package config

type Config struct {
	AdminID int64

	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string

	TelegramToken string

	LimitFree    int
	LimitPremium int

	PublicLinkMinUsers int
}
