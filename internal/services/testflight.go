package services

import (
	"context"
	"github.com/GoBotApiOfficial/gobotapi/types"
	"github.com/Laky-64/TestFlightTrackBot/internal/db"
	"github.com/Laky-64/TestFlightTrackBot/internal/db/models"
	"github.com/Laky-64/TestFlightTrackBot/internal/telegram/bot"
	"github.com/Laky-64/TestFlightTrackBot/internal/telegram/core"
	"github.com/Laky-64/TestFlightTrackBot/internal/testflight"
	"github.com/Laky-64/TestFlightTrackBot/internal/translator"
	"github.com/Laky-64/TestFlightTrackBot/internal/utils"
	"github.com/Laky-64/gologging"
	"sync"
	"time"
)

func startTestflight(ctx context.Context, rateLimit *utils.RateLimiter, b *bot.Bot, dbCtx *db.DB, tfClient *testflight.Client) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			usedLinks, err := dbCtx.LinkStore.FindUsedLinks()
			if err != nil {
				gologging.ErrorF("cleanup find: %v", err)
				continue
			}
			checkedLinks, err := tfClient.Check(usedLinks)
			if err != nil {
				gologging.ErrorF("check links: %v", err)
				continue
			}

			var updates []models.LinkUpdate
			newApps := make(map[string]*models.AppUpsert)
			var requestNotifications []models.NotificationRequest
			var removeLinks []uint
			for _, link := range usedLinks {
				checked := checkedLinks[link.URL]
				if checked.Error != nil {
					continue
				}
				if checked.Status == models.StatusInvalid {
					removeLinks = append(removeLinks, link.ID)
					continue
				}
				if _, ok := newApps[checked.AppName]; (link.AppID != nil || !ok) && len(checked.AppName) > 0 {
					newApps[checked.AppName] = &models.AppUpsert{
						LinkID:      &link.ID,
						IconURL:     checked.IconURL,
						Description: checked.Description,
					}
				}
				if checked.Status != link.Status {
					updates = append(updates, models.LinkUpdate{
						ID:      link.ID,
						AppName: checked.AppName,
						Status:  checked.Status,
					})
					requestNotifications = append(requestNotifications, models.NotificationRequest{
						LinkID: link.ID,
						Status: checked.Status,
					})
				}
			}

			err = dbCtx.AppStore.BulkUpsert(newApps)
			if err != nil {
				gologging.ErrorF("bulk upsert apps: %v", err)
				continue
			}

			if len(updates) > 0 {
				if dbCtx.LinkStore.BulkUpdate(updates) != nil {
					gologging.ErrorF("bulk update links: %v", err)
					continue
				}
				notifications, err := dbCtx.ChatLinkStore.BulkUpdateNotifications(requestNotifications)
				if err != nil {
					gologging.ErrorF("bulk update notifications: %v", err)
					continue
				}
				var wg sync.WaitGroup
				for _, n := range notifications {
					wg.Add(1)
					rateLimit.Enqueue(func() {
						defer wg.Done()
						updateContext := core.NewLightContext(b.Api, n.Lang)
						var messageKey translator.Key
						var keyboard *types.InlineKeyboardMarkup
						if n.Status == models.StatusAvailable {
							messageKey = translator.BetaOpened
							keyboard = &types.InlineKeyboardMarkup{
								InlineKeyboard: [][]types.InlineKeyboardButton{
									{
										{
											Text: updateContext.Translator.T(translator.JoinBetaBtn),
											URL:  n.LinkURL,
										},
									},
								},
							}
						} else {
							messageKey = translator.BetaClosed
						}
						_ = updateContext.SendMessageWithKeyboard(
							n.ChatID,
							updateContext.Translator.TWithData(
								messageKey,
								map[string]string{
									"AppName": n.AppName,
								},
							),
							keyboard,
						)
					})
				}
				wg.Wait()
			}
			if len(removeLinks) > 0 {
				notifications, err := dbCtx.LinkStore.BulkDelete(removeLinks)
				if err != nil {
					gologging.ErrorF("bulk delete links: %v", err)
					continue
				}
				var wg sync.WaitGroup
				for _, n := range notifications {
					wg.Add(1)
					rateLimit.Enqueue(func() {
						defer wg.Done()
						updateContext := core.NewLightContext(b.Api, n.Lang)
						_ = updateContext.SendMessage(
							n.ChatID,
							updateContext.Translator.TWithData(
								translator.TrackingRemovedError,
								map[string]string{
									"LinkURL": n.LinkURL,
								},
							),
						)
					})
				}
			}
		}
	}
}
