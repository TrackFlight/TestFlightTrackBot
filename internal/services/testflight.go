package services

import (
	"context"
	"errors"
	"github.com/GoBotApiOfficial/gobotapi/types"
	"github.com/Laky-64/gologging"
	"github.com/TrackFlight/TestFlightTrackBot/internal/config"
	"github.com/TrackFlight/TestFlightTrackBot/internal/db"
	"github.com/TrackFlight/TestFlightTrackBot/internal/db/models"
	"github.com/TrackFlight/TestFlightTrackBot/internal/telegram/bot"
	"github.com/TrackFlight/TestFlightTrackBot/internal/telegram/core"
	"github.com/TrackFlight/TestFlightTrackBot/internal/testflight"
	"github.com/TrackFlight/TestFlightTrackBot/internal/translator"
	"github.com/TrackFlight/TestFlightTrackBot/internal/utils"
	"github.com/jackc/pgx/v5"
	"maps"
	"slices"
	"sync"
	"time"
)

func startTestflight(ctx context.Context, rateLimit *utils.RateLimiter, b *bot.Bot, cfg *config.Config, dbCtx *db.DB, tfClient *testflight.Client) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			usedLinks, err := dbCtx.LinkStore.GetUsedLinks(cfg.PublicLinkMinUsers)
			if err != nil {
				gologging.ErrorF("cleanup find: %v", err)
				continue
			}
			checkedLinks, err := tfClient.Check(usedLinks)
			if err != nil {
				gologging.ErrorF("check links: %v", err)
				continue
			}

			var updates []db.BulkUpdateLinkParams
			newApps := make(map[string]db.BulkUpsertAppParams)
			var requestNotifications []db.BulkUpdateNotificationsChatParams
			var removeLinks []int64
			for _, link := range usedLinks {
				checked := checkedLinks[link.URL]
				newPublicStatus := link.IsPublic || checked.IsPublic
				if checked.Error != nil {
					continue
				}
				if checked.Status == models.LinkStatusEnumInvalid {
					removeLinks = append(removeLinks, link.ID)
					continue
				}
				if _, ok := newApps[checked.AppName]; (link.AppID.Valid || !ok) && len(checked.AppName) > 0 {
					newApps[checked.AppName] = db.BulkUpsertAppParams{
						AppID:       link.AppID.Int64,
						AppName:     checked.AppName,
						IconURL:     checked.IconURL,
						Description: checked.Description,
					}
				}
				if checked.Status != link.Status || link.IsPublic != newPublicStatus || (!link.AppID.Valid && link.AppID.Valid != (len(checked.AppName) > 0)) {
					updates = append(updates, db.BulkUpdateLinkParams{
						LinkID:   link.ID,
						AppName:  checked.AppName,
						Status:   checked.Status,
						IsPublic: newPublicStatus,
					})
				}
				if checked.Status != link.Status {
					requestNotifications = append(requestNotifications, db.BulkUpdateNotificationsChatParams{
						LinkID: link.ID,
						Status: checked.Status,
					})
				}
			}

			err = dbCtx.AppStore.BulkUpsert(slices.Collect(maps.Values(newApps)))
			if err != nil && !errors.Is(err, pgx.ErrNoRows) {
				gologging.ErrorF("bulk upsert apps: %v", err)
				continue
			}

			if len(updates) > 0 {
				if err = dbCtx.LinkStore.BulkUpdate(updates); err != nil && !errors.Is(err, pgx.ErrNoRows) {
					gologging.ErrorF("bulk update links: %v", err)
					continue
				}
				notifications, err := dbCtx.ChatStore.BulkUpdateNotifications(
					requestNotifications,
					cfg.LimitFree,
				)
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
						if n.Status == models.LinkStatusEnumAvailable {
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
