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
    chat_links.created_at AS added_at,
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
    RETURNING id, app_id, url, status, is_public, current_version, last_availability, updated_at
),
link_row AS (
    SELECT id, app_id, url, status, is_public, current_version, last_availability, updated_at
    FROM ins
    UNION ALL
    SELECT l.id, l.app_id, l.url, l.status, l.is_public, l.current_version, l.last_availability, l.updated_at
    FROM links l
    JOIN normalized n ON l.url = n.url
),
tracking AS (
    SELECT COUNT(*) AS links_count
    FROM chat_links cl
    WHERE cl.chat_id = @chat_id
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
    INSERT INTO chat_links (chat_id, link_id, notify_available, notify_closed, last_notified_version)
    SELECT
        @chat_id,
        link_row.id,
        (@notify_available::boolean AND links_count < @free_limit::bigint),
        (@notify_closed::boolean AND links_count < @free_limit::bigint),
        link_row.current_version
    FROM link_row, tracking, limit_check
    RETURNING link_id, updated_at, created_at
)
SELECT
    ic.link_id AS id,
    COALESCE(lr.app_id, -ic.link_id)::bigint AS entity_id,
    lr.url::text AS link_url,
    lr.status,
    lr.last_availability,
    lr.is_public AS is_public,
    lr.updated_at AS last_update,
    ic.created_at AS added_at
FROM ins_cl ic
JOIN link_row lr ON lr.id = ic.link_id;


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
        UNNEST(@statuses::link_status_enum[]) AS status,
        UNNEST(@versions::bigint[]) AS version
),
ranked_links AS (
    SELECT
        ROW_NUMBER() OVER (PARTITION BY cl.chat_id ORDER BY cl.created_at) AS rn,
        cl.link_id,
        cl.chat_id,
        cl.last_notified_version,
        cl.notify_available,
        cl.notify_closed
    FROM chat_links cl
    WHERE cl.notify_available OR cl.notify_closed
),
candidates AS (
    SELECT
        rl.chat_id,
        rl.link_id,
        i.version,
        i.status
    FROM input_data i
    JOIN ranked_links rl ON rl.link_id = i.link_id
    WHERE rl.last_notified_version IS DISTINCT FROM i.version
    AND (
        rl.notify_available = (i.status = 'available')
        OR rl.notify_closed = (i.status != 'available')
    )
    AND (
        EXISTS (
            SELECT 1 FROM premium_users pu WHERE pu.chat_id = rl.chat_id
        )
        OR rl.rn <= @free_limit::bigint
    )
),
notified AS (
    UPDATE chat_links cl
    SET last_notified_version = fc.version, updated_at = NOW()
    FROM candidates fc
    WHERE cl.link_id = fc.link_id AND cl.chat_id = fc.chat_id
    RETURNING cl.chat_id, cl.link_id, fc.status
)
SELECT
    c.id AS chat_id,
    c.lang,
    a.app_name,
    l.url AS link_url,
    notified.status::link_status_enum AS status
FROM notified
JOIN links l ON l.id = notified.link_id
JOIN chats c ON c.id = notified.chat_id
JOIN apps a ON a.id = l.app_id;