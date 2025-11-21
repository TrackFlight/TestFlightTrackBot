package services

import (
	"fmt"
	"maps"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/GoBotApiOfficial/gobotapi/types"
	rawTypes "github.com/GoBotApiOfficial/gobotapi/types/raw"
	"github.com/Laky-64/gologging"
	"github.com/TrackFlight/TestFlightTrackBot/internal/config"
	"github.com/TrackFlight/TestFlightTrackBot/internal/db"
	"github.com/TrackFlight/TestFlightTrackBot/internal/db/models"
	"github.com/TrackFlight/TestFlightTrackBot/internal/graphics"
	"github.com/TrackFlight/TestFlightTrackBot/internal/telegram"
	"github.com/TrackFlight/TestFlightTrackBot/internal/telegram/bot"
	"github.com/TrackFlight/TestFlightTrackBot/internal/telegram/core"
	"github.com/TrackFlight/TestFlightTrackBot/internal/testflight"
	"github.com/TrackFlight/TestFlightTrackBot/internal/tor"
	"github.com/TrackFlight/TestFlightTrackBot/internal/translator"
	"github.com/TrackFlight/TestFlightTrackBot/internal/utils"
	"github.com/robfig/cron/v3"
)

func startWeeklyHighLights(c *cron.Cron, rateLimit *utils.RateLimiter, b *bot.Bot, cfg *config.Config, dbCtx *db.DB, torClient *tor.Client) error {
	_, err := c.AddFunc("0 14 * * SUN", func() {
		events := map[graphics.BannerType]func() ([]db.WeeklyTrendingApp, error){
			graphics.HiddenGems: dbCtx.AppStore.GetWeeklyHiddenGems,
			graphics.Reopened:   dbCtx.AppStore.GetWeeklyOpened,
			graphics.Rising:     dbCtx.AppStore.GetWeeklyRising,
			graphics.Top5:       dbCtx.AppStore.GetWeeklyTrending,
		}
		sundaysCount := int((time.Now().Unix() - ReferenceSundayTime) / (7 * 24 * 3600))
		eventsOrdered := slices.Collect(maps.Keys(events))
		slices.Sort(eventsOrdered)
		eventName := eventsOrdered[sundaysCount%len(events)]
		apps, err := events[eventName]()
		if err != nil {
			gologging.ErrorF("weekly highlights %s: %v", eventName, err)
			return
		}

		var appIds []int64
		weeklyFollowers := make(map[int64]int64)
		for _, app := range apps {
			appIds = append(appIds, app.ID)
			weeklyFollowers[app.ID] = app.WeeklyFollowers
		}
		info, err := dbCtx.AppStore.GetAppsInfo(appIds)
		if err != nil {
			gologging.ErrorF("weekly highlights %s: %v", eventName, err)
			return
		}

		transaction, err := torClient.NewTransaction(len(appIds))
		if err != nil {
			gologging.ErrorF("weekly highlights %s: %v", eventName, err)
			return
		}

		var mu sync.Mutex
		var wg sync.WaitGroup
		var errors []error
		images := make(map[string][]byte)
		pool := utils.NewPool(tor.MaxRequestsPerInstance * torClient.InstanceCount())

		for _, app := range info {
			wg.Add(1)
			pool.Enqueue(func() {
				defer wg.Done()
				var errReturn error
				iconURL := app.IconURL.String
				for i := 0; i < 3; i++ {
					request, err := transaction.ExecuteRequest(
						testflight.ChangeLinkResolution(iconURL, 480),
					)
					if err != nil {
						errReturn = err
						continue
					}
					if request.StatusCode != 200 {
						continue
					}
					errReturn = nil
					mu.Lock()
					images[iconURL] = request.Body
					mu.Unlock()
					break
				}
				if errReturn != nil {
					mu.Lock()
					errors = append(errors, fmt.Errorf("failed to download %s: %w", iconURL, errReturn))
					mu.Unlock()
				}
			})
		}
		wg.Wait()
		transaction.Close()

		if len(errors) > 0 {
			gologging.Error(errors[0])
			return
		}

		var orderedImages [][]byte
		for _, app := range info {
			orderedImages = append(orderedImages, images[app.IconURL.String])
		}

		users, err := dbCtx.PreferencesStore.GetAllNotifiableWeeklyInsightUsers()
		if err != nil {
			gologging.ErrorF("weekly highlights get users: %v", err)
			return
		}

		languageBanners := make(map[string]rawTypes.InputFile)

		translatorEventName := strings.ToUpper(string(eventName))
		eventTitleKey := translator.Key(fmt.Sprintf("EVENT_%s_TITLE", translatorEventName))
		for _, user := range users {
			t := translator.New(user.Lang)
			if _, ok := languageBanners[t.Language()]; !ok {
				banner, errGen := graphics.GenerateBanner(
					eventName,
					t.T(eventTitleKey),
					t.T(translator.EventSubtitle),
					orderedImages,
				)
				if errGen != nil {
					gologging.ErrorF("failed to generate banner: %v", errGen)
					return
				}
				languageBanners[t.Language()] = types.InputBytes{
					Name: fmt.Sprintf("banner_%s.png", t.Language()),
					Data: banner,
				}
			}
		}

		var wgMessages sync.WaitGroup
		var unreachableChats []db.BulkUpdateNotifiableUsersChatParams
		mutexByLang := make(map[string]*sync.Mutex)
		for _, user := range users {
			wg.Add(1)
			rateLimit.Enqueue(func() {
				defer wg.Done()
				ctx := core.NewLightContext(b.Api, user.Lang)
				lang := ctx.Translator.Language()
				var sBuilder strings.Builder
				sBuilder.WriteString(ctx.Translator.TWithData(
					translator.Key(fmt.Sprintf("EVENT_%s_TEXT_TITLE", translatorEventName)),
					map[string]string{
						"Title": ctx.Translator.T(eventTitleKey),
					},
				))
				sBuilder.WriteRune('\n')
				sBuilder.WriteString(ctx.Translator.T(translator.Key(fmt.Sprintf("EVENT_%s_TEXT_DESCRIPTION", translatorEventName))))
				sBuilder.WriteString("\n\n")
				for i, weekAppInfo := range apps {
					appInfo := info[i]
					var tempSBuilder strings.Builder
					_ = weekAppInfo
					tempSBuilder.WriteString("<i>")
					tempSBuilder.WriteString(
						ctx.Translator.TWithDataCountable(
							"APP_FOLLOWERS_AMOUNT",
							map[string]string{
								"Amount": strconv.Itoa(int(appInfo.Followers)),
							},
							int(appInfo.Followers),
						),
					)
					if weekAppInfo.WeeklyFollowers != 0 {
						tempSBuilder.WriteString(fmt.Sprintf(" (+%d)", weekAppInfo.WeeklyFollowers))
					}
					tempSBuilder.WriteString("</i>")
					tempSBuilder.WriteString(" • ")
					var availability string
					switch weekAppInfo.Status {
					case models.LinkStatusEnumAvailable:
						availability = ctx.Translator.T(
							translator.StatusAvailable,
						)
					case models.LinkStatusEnumFull:
						availability = ctx.Translator.T(
							translator.StatusFull,
						)
					case models.LinkStatusEnumClosed:
						availability = ctx.Translator.T(
							translator.StatusClosed,
						)
					}
					tempSBuilder.WriteString(availability)

					sBuilder.WriteString(fmt.Sprintf(
						"%d. <b>%s</b>\n%s\n\n",
						i+1,
						appInfo.AppName.String,
						tempSBuilder.String(),
					))
				}
				mu.Lock()
				if _, ok := mutexByLang[lang]; !ok {
					mutexByLang[lang] = &sync.Mutex{}
				}
				m := mutexByLang[lang]
				mu.Unlock()
				m.Lock()
				banner := languageBanners[lang]

				if _, ok := banner.(types.InputURL); ok {
					m.Unlock()
				}
				update, err2 := ctx.SendPhotoWithKeyboard(
					user.ID,
					sBuilder.String(),
					banner,
					&types.InlineKeyboardMarkup{
						InlineKeyboard: [][]types.InlineKeyboardButton{
							{
								{
									Text: ctx.Translator.T(translator.OpenMiniappBtn),
									WebApp: &types.WebAppInfo{
										URL: cfg.MiniAppURL,
									},
								},
							},
						},
					},
				)
				if err2 == nil {
					if _, ok := banner.(types.InputBytes); ok {
						var bestImage types.PhotoSize
						photoSizes := update.Result.(types.Message).Photo
						for _, img := range photoSizes {
							if img.FileSize > bestImage.FileSize {
								bestImage = img
							}
						}
						mu.Lock()
						languageBanners[lang] = types.InputURL(bestImage.FileID)
						mu.Unlock()
					}
				} else {
					mu.Lock()
					if status := telegram.MapErrorToUserStatus(err2); status != models.UserStatusEnumReachable {
						unreachableChats = append(unreachableChats, db.BulkUpdateNotifiableUsersChatParams{
							ChatID: user.ID,
							Status: status,
						})
					}
					mu.Unlock()
				}

				if _, ok := banner.(types.InputBytes); ok {
					m.Unlock()
				}
			})
		}
		wgMessages.Wait()

		if len(unreachableChats) > 0 {
			err = dbCtx.ChatStore.BulkUpdateNotifiableUsers(unreachableChats)
			if err != nil {
				gologging.ErrorF("bulk update unreachable chats: %v", err)
				return
			}
		}
	})
	return err
}
