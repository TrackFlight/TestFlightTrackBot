package db

import (
	"github.com/Laky-64/TestFlightTrackBot/internal/config"
	"github.com/Laky-64/TestFlightTrackBot/internal/db/models"
	"gorm.io/gorm"
)

type linkStore struct {
	db  *gorm.DB
	cfg *config.Config
}

func (ctx *linkStore) Add(url string) (uint, error) {
	link := models.Link{
		URL: url,
	}
	err := ctx.db.FirstOrCreate(&link, `url = ?`, url).Error
	return link.ID, err
}

func (ctx *linkStore) FindByURL(url string) (models.Link, error) {
	var link models.Link
	err := ctx.db.First(&link, "url = ?", url).Error
	return link, err
}

func (ctx *linkStore) FindUsedLinks() ([]models.Link, error) {
	var links []models.Link
	err := ctx.db.
		Table("links").
		Select("DISTINCT ON (links.id) links.*").
		Joins("JOIN chat_links ON chat_links.link_id = links.id").
		Scan(&links).Error
	return links, err
}

func (ctx *linkStore) BulkUpdate(updates []models.LinkUpdate) error {
	var ids []uint
	var appNames []string
	var statuses []models.LinkStatus
	for _, u := range updates {
		ids = append(ids, u.ID)
		appNames = append(appNames, u.AppName)
		statuses = append(statuses, u.Status)
	}
	return bulkExec(
		ctx.db,
		`
			WITH input_data AS (
				SELECT *
				FROM UNNEST(
					ARRAY[?]::bigint[],
					ARRAY[?]::link_status_enum[],
					ARRAY[?]::text[]
				) AS i(link_id, status, app_name)
			)
			UPDATE links
			SET status = i.status,
				updated_at = NOW(),
				app_id = COALESCE(links.app_id, apps.id),
				last_availability = CASE 
				WHEN i.status = 'available' THEN NOW()
					ELSE links.last_availability
				END
			FROM input_data i
			LEFT JOIN apps ON apps.app_name = i.app_name
			WHERE links.id = i.link_id
			AND (
				links.status IS DISTINCT FROM i.status
				OR links.app_id IS DISTINCT FROM apps.id
			);
		`,
		ids, statuses, appNames,
	)
}

func (ctx *linkStore) BulkDelete(ids []uint) ([]models.NotifiedDeleteChat, error) {
	var notifiedDeleteChats []models.NotifiedDeleteChat
	err := bulkRaw(
		ctx.db,
		`
			WITH
			input_data AS (
				SELECT UNNEST(
				    ARRAY[?]::bigint[]
				)
			),
			deleted AS (
				DELETE FROM links
				USING input_data i
				WHERE links.id = i.unnest
				RETURNING links.id AS link_id, links.url
			)
			SELECT chats.id AS chat_id, chats.lang, deleted.link_id, deleted.url AS link_url
			FROM deleted
			JOIN chat_links ON chat_links.link_id = deleted.link_id
			JOIN chats ON chats.id = chat_links.chat_id;
		`,
		&notifiedDeleteChats,
		ids,
	)
	return notifiedDeleteChats, err
}
