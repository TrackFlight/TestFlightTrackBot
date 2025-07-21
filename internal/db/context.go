package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Laky-64/TestFlightTrackBot/internal/config"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/valkey-io/valkey-go"
	"log"
)

func NewDB(cfg *config.Config) (*DB, error) {
	dsn := fmt.Sprintf(
		"host=db port=5432 user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
	)
	if err := goose.SetDialect("pgx"); err != nil {
		log.Fatalf("goose: failed to set dialect: %v", err)
	}

	sqlConn, err := goose.OpenDBWithDriver("pgx", dsn)
	if err != nil {
		log.Fatalf("goose: failed to open DB: %v", err)
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(sqlConn)

	goose.SetLogger(goose.NopLogger())

	if err = goose.Up(sqlConn, "internal/db/schema"); err != nil {
		return nil, fmt.Errorf("goose up: %w", err)
	}

	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	redis, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: []string{"valkey:6379"},
	})
	if err != nil {
		return nil, fmt.Errorf("connect to valkey: %w", err)
	}
	return new(conn, redis), nil
}
