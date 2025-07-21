-- name: Search :many
-- cache: type:get table:apps key:name
WITH apps_fuzzy AS (
    SELECT
        id,
        app_name,
        levenshtein(app_name, @name::text)::float / length(app_name) AS l_value
    FROM apps
    WHERE levenshtein(app_name, @name::text)::float / length(app_name) <= 1.4
    LIMIT 3
)
SELECT
    a.id AS app_id,
    a.app_name,
    la.links_count,
    la.followers,
    la.updated_at::timestamptz AS updated_at
FROM apps_fuzzy a
CROSS JOIN LATERAL (
    SELECT
        COUNT(l.*)        AS links_count,
        COUNT(cl.*)       AS followers,
        MAX(l.updated_at) AS updated_at
    FROM links l
    LEFT JOIN chat_links cl ON cl.link_id = l.id
    WHERE l.app_id = a.id AND l.status IS NOT NULL
) la
ORDER BY a.l_value, la.followers DESC;


-- name: BulkUpsert :exec
WITH input_data AS (
    SELECT
        UNNEST(@app_ids::bigint[])    AS app_id,
        UNNEST(@app_names::text[])    AS app_name,
        UNNEST(@icon_urls::text[])    AS icon_url,
        UNNEST(@descriptions::text[]) AS description
),
updated AS (
    UPDATE apps AS a
    SET
        app_name    = i.app_name,
        icon_url    = i.icon_url,
        description = i.description,
        updated_at  = NOW()
    FROM input_data AS i
    WHERE a.app_name = i.app_name OR a.id = i.app_id
    RETURNING a.id, a.app_name
)
INSERT INTO apps (app_name, icon_url, description)
SELECT
    i.app_name,
    i.icon_url,
    i.description
FROM input_data as i
WHERE NOT EXISTS (
    SELECT 1 FROM updated
    WHERE updated.id = i.app_id
    OR updated.app_name = i.app_name
);