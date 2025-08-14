-- name: Search :many
SELECT
    id
FROM apps
WHERE levenshtein(lower(app_name), lower(@name::text))::float / GREATEST(length(app_name) - length(@name::text), 1) <= 1.4
ORDER BY
    ((lower(app_name) LIKE lower(@name || '%'))::int * 100) DESC,
    ((lower(app_name) LIKE lower('%' || @name || '%'))::int * 50) DESC,
    levenshtein(left(lower(app_name), 3), left(lower(@name::text), 3))::float,
    levenshtein(lower(app_name), lower(@name::text))::float / GREATEST(length(app_name) - length(@name), 1)
LIMIT 5;


-- name: GetTrending :many
-- cache: type:get table:trending_apps ttl:1h
SELECT
    a.id
FROM apps a
JOIN links l ON l.app_id = a.id
JOIN chat_links cl ON cl.link_id = l.id
GROUP BY a.id
ORDER BY
    COUNT(DISTINCT cl.chat_id) * 0.2 +
    COUNT(DISTINCT cl.chat_id) FILTER (
        WHERE cl.created_at >= NOW() - INTERVAL '7 days'
    ) * 0.8
    DESC
LIMIT 50;


-- name: GetAppsInfo :many
-- cache: type:get table:apps key:entity_ids ttl:30m
WITH
app_ids AS (
    SELECT ABS(x) AS id FROM unnest(@entity_ids::bigint[]) x WHERE x >= 0
),
link_ids AS (
    SELECT ABS(x) AS id FROM unnest(@entity_ids::bigint[]) x WHERE x < 0
),
link_follower_counts AS (
    SELECT
        links.id,
        COUNT(DISTINCT cl.chat_id) AS followers
    FROM links
    JOIN chat_links cl ON cl.link_id = links.id
    JOIN link_ids ON link_ids.id = links.id
    GROUP BY links.id
),
app_follower_counts AS (
    SELECT
        apps.id,
        apps.app_name,
        apps.icon_url,
        apps.description,
        (
            SELECT COUNT(DISTINCT cl.chat_id)
            FROM links l
            JOIN chat_links cl ON cl.link_id = l.id
            WHERE l.app_id = apps.id
        ) AS followers
    FROM apps
    JOIN app_ids ON app_ids.id = apps.id
    GROUP BY apps.id
)
SELECT
    -lf.id AS entity_id,
    NULL::varchar AS app_name,
    NULL::varchar AS icon_url,
    NULL::varchar AS description,
    lf.followers
FROM link_follower_counts lf
UNION ALL
SELECT
    af.id AS entity_id,
    af.app_name,
    af.icon_url,
    af.description,
    af.followers
FROM app_follower_counts af;


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
        icon_url    = NULLIF(i.icon_url, ''),
        description = NULLIF(i.description, ''),
        updated_at  = NOW()
    FROM input_data AS i
    WHERE a.app_name = i.app_name OR a.id = i.app_id
    RETURNING a.id, a.app_name
)
INSERT INTO apps (app_name, icon_url, description)
SELECT
    i.app_name,
    NULLIF(i.icon_url, ''),
    NULLIF(i.description, '')
FROM input_data as i
WHERE NOT EXISTS (
    SELECT 1 FROM updated
    WHERE updated.id = i.app_id
    OR updated.app_name = i.app_name
);