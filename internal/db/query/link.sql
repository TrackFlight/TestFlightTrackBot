-- name: GetUsedLinks :many
SELECT
    links.id,
    links.url,
    links.app_id,
    links.status,
    CASE
    WHEN COUNT(chat_links.chat_id) > @min_public::bigint THEN TRUE
    WHEN links.is_public THEN TRUE
    ELSE FALSE
    END AS is_public
FROM links
JOIN chat_links ON chat_links.link_id = links.id
GROUP BY links.id, links.url, links.app_id, links.status;


-- name: BulkDelete :many
WITH deleted AS (
    DELETE FROM links
    WHERE links.id = ANY(@link_ids::bigint[])
    RETURNING links.id AS link_id, links.url
)
SELECT chats.id AS chat_id, chats.lang, deleted.link_id, deleted.url AS link_url
FROM deleted
JOIN chat_links ON chat_links.link_id = deleted.link_id
JOIN chats ON chats.id = chat_links.chat_id;


-- name: GetLinksByApps :many
-- cache: type:get-many table:app_links key:app_id ttl:1d
SELECT
    links.id,
    (
        CASE
        WHEN links.is_public THEN links.url
        ELSE
            'testflight.apple.com/join/' ||
            repeat(
                'X',
                length(substring(links.url from 'testflight\.apple\.com/join/(.+)$'))
            )
        END
    )::text AS url,
    links.app_id::bigint AS app_id,
    links.status,
    links.last_availability,
    links.is_public,
    links.updated_at AS last_update
FROM links
WHERE links.app_id = ANY(@app_id::bigint[])
ORDER BY links.created_at;


-- name: BulkUpdate :many
-- cache: type:update_version table:links key:link_ids ttl:1w
-- exclude: link_id
WITH input_data AS (
    SELECT
        UNNEST(@link_ids::bigint[]) AS link_id,
        UNNEST(@app_names::text[]) AS app_name,
        UNNEST(@statuses::link_status_enum[]) AS status,
        UNNEST(@is_public::boolean[]) AS is_public
)
UPDATE links
SET status = i.status,
    updated_at = NOW(),
    app_id = COALESCE(links.app_id, apps.id),
    last_availability = CASE
    WHEN i.status = 'available' THEN NOW()
        ELSE links.last_availability
    END,
    is_public = i.is_public
FROM input_data i
LEFT JOIN apps ON apps.app_name = i.app_name
WHERE links.id = i.link_id
AND (
    links.status IS DISTINCT FROM i.status
    OR links.app_id IS DISTINCT FROM apps.id
    OR links.is_public IS DISTINCT FROM i.is_public
) RETURNING links.id AS link_id;