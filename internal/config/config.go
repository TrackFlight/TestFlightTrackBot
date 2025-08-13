package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

func Load() (*Config, error) {
	cfg := &Config{
		AdminID:       int64(getEnvInt("ADMIN_USER_ID", 0)),
		DBUser:        os.Getenv("DB_USER"),
		DBPassword:    os.Getenv("DB_PASSWORD"),
		DBName:        os.Getenv("DB_NAME"),
		TelegramToken: os.Getenv("TELEGRAM_TOKEN"),

		LimitFree:          int64(getEnvInt("LIMIT_FREE", 3)),
		LimitPremium:       int64(getEnvInt("LIMIT_PREMIUM", 10)),
		MaxFollowingLinks:  int64(getEnvInt("MAX_FOLLOWING_LINKS", 50)),
		PublicLinkMinUsers: int64(getEnvInt("PUBLIC_LINK_MIN_USERS", 20)),

		MiniAppURL: fmt.Sprintf("https://%s", os.Getenv("SERVER_NAME")),
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
