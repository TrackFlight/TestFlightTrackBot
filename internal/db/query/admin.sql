-- name: GetBotStats :one
WITH
    chat_stats AS (
        SELECT COUNT(*) AS total_chats FROM chats
    ),
    chat_link_stats AS (
        SELECT
            COUNT(DISTINCT chat_id) AS active_chats,
            COUNT(DISTINCT link_id) AS tracked_links
        FROM chat_links
    ),
    link_stats AS (
        SELECT COUNT(*) AS total_links FROM links
    )
SELECT
    cls.active_chats,
    cs.total_chats,
    cls.tracked_links,
    ls.total_links
FROM chat_stats cs, chat_link_stats cls, link_stats ls;