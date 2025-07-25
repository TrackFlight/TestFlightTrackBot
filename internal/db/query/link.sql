-- name: GetUsedLinks :many
SELECT DISTINCT ON (links.id) links.id, links.url, links.app_id, links.status
FROM links
JOIN chat_links ON chat_links.link_id = links.id;

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


-- name: BulkUpdate :many
-- cache: type:update_version table:links key:link_ids ttl:1w
-- exclude: link_id
WITH input_data AS (
    SELECT
        UNNEST(@link_ids::bigint[]) AS link_id,
        UNNEST(@app_names::text[]) AS app_name,
        UNNEST(@statuses::link_status_enum[]) AS status
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
) RETURNING links.id AS link_id;