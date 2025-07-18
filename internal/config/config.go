package config

import (
	"errors"
	"os"
	"strconv"
)

func Load() (*Config, error) {
	cfg := &Config{
		AdminID:       int64(getEnvInt("ADMIN_USER_ID", 0)),
		DBHost:        os.Getenv("DB_HOST"),
		DBUser:        os.Getenv("DB_USER"),
		DBPassword:    os.Getenv("DB_PASSWORD"),
		DBName:        os.Getenv("DB_NAME"),
		TelegramToken: os.Getenv("TELEGRAM_TOKEN"),

		LimitFree:          getEnvInt("LIMIT_FREE", 3),
		LimitPremium:       getEnvInt("LIMIT_PREMIUM", 10),
		PublicLinkMinUsers: getEnvInt("PUBLIC_LINK_MIN_USERS", 20),
	}

	if len(cfg.DBHost) == 0 {
		return nil, errors.New("DB_HOST environment variable not set")
	}
	if len(cfg.DBUser) == 0 {
		return nil, errors.New("DB_USER environment variable not set")
	}
	if len(cfg.DBPassword) == 0 {
		return nil, errors.New("DB_PASSWORD environment variable not set")
	}
	if len(cfg.DBName) == 0 {
		return nil, errors.New("DB_NAME environment variable not set")
	}
	if len(cfg.TelegramToken) == 0 {
		return nil, errors.New("TELEGRAM_TOKEN environment variable not set")
	}
	if cfg.AdminID == 0 {
		return nil, errors.New("ADMIN_USER_ID environment variable not set or invalid")
	}
	return cfg, nil
}

func getEnvInt(key string, fallback int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return fallback
}
