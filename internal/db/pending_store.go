package db

import (
	"github.com/Laky-64/TestFlightTrackBot/internal/db/models"
	"gorm.io/gorm"
	"sync"
)

type pendingStore struct {
	db              *gorm.DB
	pendingTracking sync.Map
}

func (ctx *pendingStore) Add(chatID int64) error {
	if _, exists := ctx.pendingTracking.Swap(chatID, true); exists {
		return nil
	}

	var pendingTrack models.PendingTrack
	result := ctx.db.FirstOrCreate(&pendingTrack, &models.PendingTrack{
		ChatID: chatID,
	})
	return result.Error
}

func (ctx *pendingStore) LoadFromDB() error {
	var pendingTracks []models.PendingTrack
	if err := ctx.db.Find(&pendingTracks).Error; err != nil {
		return err
	}

	for _, track := range pendingTracks {
		ctx.pendingTracking.Store(track.ChatID, true)
	}
	return nil
}

func (ctx *pendingStore) Exists(chatID int64) bool {
	_, exists := ctx.pendingTracking.Load(chatID)
	return exists
}

func (ctx *pendingStore) Remove(chatID int64) bool {
	if _, exists := ctx.pendingTracking.LoadAndDelete(chatID); !exists {
		return false
	}
	ctx.db.Delete(models.PendingTrack{
		ChatID: chatID,
	})
	return true
}
