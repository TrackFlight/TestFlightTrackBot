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
    links.url AS link_url,
    COALESCE(links.app_id, -links.id)::bigint AS entity_id,
    links.status,
    links.is_public,
    chat_links.notify_available,
    chat_links.notify_closed,
    links.last_availability,
    links.updated_at AS last_update
FROM chat_links
JOIN links ON chat_links.link_id = links.id
WHERE chat_id = @chat_id
ORDER BY chat_links.created_at;


-- name: Track :one
-- cache: type:remove table:chat_links key:chat_id fields:all_by_key
-- order: chat_id, link_url, notify_available, notify_closed, free_limit, max_following_links
WITH normalized AS (
    SELECT REGEXP_REPLACE(@link_url, '^https?://', '') AS url
),
ins AS (
    INSERT INTO links (url)
    SELECT url FROM normalized n
    WHERE NOT EXISTS (
        SELECT 1 FROM links l WHERE l.url = n.url
    )
    RETURNING id, url, updated_at
),
link_row AS (
    SELECT id, url, updated_at
    FROM ins
    UNION ALL
    SELECT l.id, l.url, l.updated_at
    FROM links l
    JOIN normalized n ON l.url = n.url
),
tracking AS (
    SELECT COUNT(*) FILTER (
        WHERE cl.chat_id = @chat_id
        AND (cl.notify_available OR cl.notify_closed)
    )
    AS links_count
    FROM chat_links cl
),
limit_check AS (
    SELECT assert(
        links_count < @max_following_links::bigint,
        format('Too many links followed: %s / %s', links_count, @max_following_links),
        'P1600'
    )
    FROM tracking
),
ins_cl AS (
    INSERT INTO chat_links (chat_id, link_id, notify_available, notify_closed)
    SELECT
        @chat_id,
        link_row.id,
        (@notify_available::boolean AND links_count < @free_limit::bigint),
        (@notify_closed::boolean AND links_count < @free_limit::bigint)
    FROM link_row, tracking, limit_check
    RETURNING link_id
)
SELECT
    ic.link_id AS id,
    COALESCE(el.app_id, -ic.link_id)::bigint AS entity_id,
    lr.url::text AS link_url,
    el.status,
    el.last_availability,
    COALESCE(el.is_public, false)::boolean AS is_public,
    lr.updated_at   AS last_update
FROM ins_cl ic
LEFT JOIN links el      ON el.id = ic.link_id
LEFT JOIN link_row lr   ON lr.id = ic.link_id;


-- name: UpdateNotificationSettings :exec
-- cache: type:remove table:chat_links key:chat_id fields:all_by_key
-- order: chat_id, link_id, notify_available, notify_closed, free_limit
WITH limit_check AS (
    SELECT assert(
        COUNT(*) < @free_limit::bigint,
        format('Too many active tracked links: %s / %s', COUNT(*), @free_limit),
        'P1600'
    )
    FROM chat_links cl
    WHERE cl.chat_id = @chat_id::bigint
    AND cl.link_id != @link_id::bigint
    AND (cl.notify_available OR cl.notify_closed)
)
UPDATE chat_links
SET notify_available = @notify_available,
notify_closed = @notify_closed,
updated_at = NOW()
FROM limit_check
WHERE chat_id = @chat_id::bigint
AND link_id = @link_id::bigint
AND (
    notify_available IS DISTINCT FROM @notify_available
    OR notify_closed IS DISTINCT FROM @notify_closed
);


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
        ROW_NUMBER() OVER (PARTITION BY cl.chat_id ORDER BY cl.created_at) AS rn
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
        cl.notify_available = (i.status = 'available')
        OR cl.notify_closed = (i.status != 'available')
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