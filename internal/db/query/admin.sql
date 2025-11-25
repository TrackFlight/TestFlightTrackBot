-- name: GetBotStats :one
WITH
    chat_stats AS (
        SELECT
            COUNT(*) FILTER (
                WHERE status = 'reachable'
            ) AS total_chats,
            COUNT(*) FILTER (
                WHERE status = 'blocked_by_user'
            ) AS blocked_chats
        FROM chats
    ),
    chat_link_stats AS (
        SELECT
            COUNT(DISTINCT chat_id) AS active_chats,
            COUNT(DISTINCT link_id) AS tracked_links
        FROM chat_links
        JOIN chats ON chats.id = chat_links.chat_id
        WHERE chats.status = 'reachable'
    ),
    link_stats AS (
        SELECT COUNT(*) AS total_links FROM links
    )
SELECT
    cls.active_chats,
    cs.total_chats,
    cs.blocked_chats,
    cls.tracked_links,
    ls.total_links
FROM chat_stats cs, chat_link_stats cls, link_stats ls;