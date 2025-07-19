package db

import (
	"github.com/Laky-64/TestFlightTrackBot/internal/config"
	"github.com/Laky-64/TestFlightTrackBot/internal/db/models"
	"github.com/Laky-64/TestFlightTrackBot/internal/translator"
	"gorm.io/gorm"
)

type chatLinkStore struct {
	db  *gorm.DB
	cfg *config.Config
}

func (ctx *chatLinkStore) GetLanguage(chatID int64, languageHint string) string {
	if !translator.IsSupported(languageHint) {
		languageHint = ""
	}
	var chat models.Chat
	result := ctx.db.First(&chat, "id = ?", chatID)
	if result.RowsAffected == 0 {
		ctx.db.Create(&models.Chat{
			ID:   chatID,
			Lang: languageHint,
		})
	}
	return chat.Lang
}

func (ctx *chatLinkStore) Track(chatID int64, linkId uint, linkURL string) (*models.TrackingLink, error) {
	var trackedLink models.TrackingLink
	err := ctx.db.Raw(`
		WITH existing_link AS (
			SELECT id, app_id, status, last_availability
			FROM links
			WHERE url = ? OR id = ?
			LIMIT 1
		),
		inserted_link AS (
			INSERT INTO links (url, created_at, updated_at)
			SELECT ?, NOW(), NOW()
			WHERE NOT EXISTS (SELECT 1 FROM existing_link)
			RETURNING id
		),
		final_link AS (
			SELECT id FROM inserted_link
			UNION ALL
			SELECT id FROM existing_link
		),
		tracking AS (
			SELECT COUNT(*) AS links_count
			FROM chat_links cl
			JOIN links l ON l.id = cl.link_id
			WHERE cl.chat_id = ?
		),
		inserted_tracking AS (
			INSERT INTO chat_links (chat_id, link_id, allow_opened, created_at, updated_at)
			VALUES (
				?,
				(SELECT id FROM final_link),
				(SELECT links_count FROM tracking) < 2,
				NOW(),
				NOW()
			)
            ON CONFLICT (chat_id, link_id) DO NOTHING
			RETURNING link_id
		)
		SELECT inserted_tracking.link_id AS id, apps.app_name, apps.icon_url, apps.description, existing_link.status, existing_link.last_availability FROM inserted_tracking
		LEFT JOIN existing_link ON existing_link.id = inserted_tracking.link_id
		LEFT JOIN apps ON apps.id = existing_link.app_id;
	`, linkURL, linkId, linkURL, chatID, chatID).Scan(&trackedLink).Error
	if err != nil {
		return nil, err
	}
	if trackedLink.ID == 0 {
		return nil, nil
	}
	return &trackedLink, nil
}

func (ctx *chatLinkStore) TrackedList(chatID int64) ([]models.TrackingLink, error) {
	var chatLinks []models.TrackingLink
	if err := ctx.db.Table("chat_links").
		Select("links.id, apps.app_name, apps.icon_url, apps.description, links.status, links.last_availability").
		Joins("JOIN links ON chat_links.link_id = links.id").
		Joins("JOIN apps ON links.app_id = apps.id").
		Where("chat_id = ?", chatID).
		Order("chat_links.created_at").
		Find(&chatLinks).Error; err != nil {
		return nil, err
	}
	return chatLinks, nil
}

func (ctx *chatLinkStore) Delete(chatID int64, linkId uint) error {
	result := ctx.db.
		Where("chat_id = ? AND link_id = ?", chatID, linkId).
		Delete(&models.ChatLink{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (ctx *chatLinkStore) TrackedCount(chatID int64) (int, error) {
	var count int64
	if err := ctx.db.Table("chat_links").
		Where("chat_id = ?", chatID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (ctx *chatLinkStore) BulkUpdateNotifications(updates []models.NotificationRequest) ([]models.NotificationChat, error) {
	var notifiedChats []models.NotificationChat
	var ids []uint
	var statuses []models.LinkStatus
	for _, u := range updates {
		ids = append(ids, u.LinkID)
		statuses = append(statuses, u.Status)
	}
	err := bulkRaw(
		ctx.db,
		`
			WITH input_data AS (
				SELECT *
				FROM UNNEST(
					ARRAY[?]::bigint[],
					ARRAY[?]::link_status_enum[]
				) AS i(link_id, status)
			),
			ranked_links AS (
				SELECT
					cl.link_id,
					cl.chat_id,
					ROW_NUMBER() OVER (PARTITION BY cl.chat_id ORDER BY l.created_at) AS rn
				FROM chat_links cl
				JOIN links l ON l.id = cl.link_id
			),
			candidates AS (
				SELECT
					cl.chat_id,
        			cl.link_id,
        			i.status
				FROM input_data i
				JOIN chat_links cl ON cl.link_id = i.link_id
				JOIN ranked_links rl ON rl.link_id = cl.link_id AND rl.chat_id = cl.chat_id
				WHERE cl.last_notified_status IS DISTINCT FROM i.status
				AND (
					cl.allow_opened = (i.status = 'available')
					OR cl.allow_closed = (i.status != 'available')
				)
				AND (
					EXISTS (
						SELECT 1 FROM premium_users pu WHERE pu.chat_id = cl.chat_id
					)
					OR rl.rn <= ?
				)
			),
			notified AS (
				UPDATE chat_links cl
				SET last_notified_status = fc.status, updated_at = NOW()
				FROM candidates fc
				WHERE cl.link_id = fc.link_id AND cl.chat_id = fc.chat_id
				RETURNING cl.chat_id, cl.link_id, cl.last_notified_status AS status
			)
			SELECT
				c.id AS chat_id,
				c.lang,
				a.app_name,
				l.url AS link_url,
				notified.status
			FROM notified
			JOIN links l ON l.id = notified.link_id
			JOIN chats c ON c.id = notified.chat_id
			JOIN apps a ON a.id = l.app_id;
		`,
		&notifiedChats,
		ids,
		statuses,
		ctx.cfg.LimitFree,
	)
	return notifiedChats, err
}
