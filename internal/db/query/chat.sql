-- name: GetLanguage :one
-- cache: type:get table:chats key:id
WITH upsert AS (
    INSERT INTO chats (id, lang)
    VALUES (@id, @lang)
    ON CONFLICT (id) DO NOTHING
    RETURNING lang
)
SELECT lang FROM upsert
UNION
SELECT lang FROM chats WHERE id = @id;


-- name: TrackedList :many
-- cache: type:get table:chat_links key:chat_id version_by:links.id
SELECT
    links.id AS id,
    REGEXP_REPLACE(links.url, '^https?://', '') AS link_url,
    COALESCE(links.app_id, -links.id)::bigint AS entity_id,
    links.status,
    links.last_availability,
    links.updated_at AS last_update,
    links.is_public
FROM chat_links
JOIN links ON chat_links.link_id = links.id
WHERE chat_id = @chat_id
ORDER BY chat_links.created_at;


-- name: Track :one
-- cache: type:remove table:chat_links key:chat_id fields:all_by_key
-- order: chat_id, link_id, link_url, free_limit
WITH existing_link AS (
    SELECT id, url, app_id, status, is_public, last_availability, updated_at
    FROM links as l
    WHERE
        l.url = @link_url
        OR l.id = @link_id
    LIMIT 1
),
inserted_link AS (
    INSERT INTO links (url)
    SELECT @link_url
    WHERE NOT EXISTS (SELECT 1 FROM existing_link)
    RETURNING id, url, updated_at
),
final_link AS (
    SELECT id, url, updated_at FROM inserted_link
    UNION ALL
    SELECT id, url, updated_at FROM existing_link
),
tracking AS (
    SELECT COUNT(*) AS links_count
    FROM chat_links cl
    JOIN links l ON l.id = cl.link_id
    WHERE cl.chat_id = @chat_id
    AND cl.allow_opened OR cl.allow_closed
),
inserted_tracking AS (
    INSERT INTO chat_links (chat_id, link_id, allow_opened)
    VALUES (
        @chat_id,
        (SELECT id FROM final_link),
        (SELECT links_count FROM tracking) < @free_limit::bigint
    )
    ON CONFLICT (chat_id, link_id) DO NOTHING
    RETURNING link_id
)
SELECT
    inserted_tracking.link_id AS id,
    COALESCE(existing_link.app_id, -inserted_tracking.link_id)::bigint AS entity_id,
    REGEXP_REPLACE(final_link.url, '^https?://', '') AS link_url,
    existing_link.status,
    existing_link.last_availability,
    COALESCE(existing_link.is_public, false)::boolean AS is_public,
    final_link.updated_at AS last_update
FROM inserted_tracking
LEFT JOIN existing_link ON existing_link.id = inserted_tracking.link_id
LEFT JOIN final_link ON final_link.id = inserted_tracking.link_id;


-- name: Delete :exec
-- cache: type:remove table:chat_links key:chat_id fields:all_by_key
DELETE FROM chat_links
WHERE chat_id = @chat_id
AND link_id = ANY(@link_ids::bigint[]);


-- name: BulkUpdateNotifications :many
WITH input_data AS (
    SELECT
        UNNEST(@link_ids::bigint[]) AS link_id,
        UNNEST(@statuses::link_status_enum[]) AS status
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
        OR rl.rn <= @free_limit::bigint
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