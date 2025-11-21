-- name: UpdateNotificationsPreferences :exec
-- cache: type:remove table:preferences key:chat_id fields:row
INSERT INTO preferences (chat_id, new_features_notifications, weekly_insights_notifications)
VALUES (@chat_id, @new_features_notifications, @weekly_insights_notifications)
ON CONFLICT (chat_id) DO
UPDATE SET new_features_notifications = EXCLUDED.new_features_notifications,
weekly_insights_notifications = EXCLUDED.weekly_insights_notifications,
updated_at = NOW();

-- name: GetNotificationsPreferences :one
-- cache: type:get table:preferences key:chat_id
WITH upsert AS (
    INSERT INTO preferences (chat_id)
    VALUES (@chat_id)
    ON CONFLICT (chat_id) DO NOTHING
    RETURNING new_features_notifications, weekly_insights_notifications
)
SELECT * FROM upsert
UNION
SELECT
    new_features_notifications,
    weekly_insights_notifications
FROM preferences WHERE chat_id = @chat_id;


-- name: GetAllNotifiableWeeklyInsightUsers :many
SELECT id, lang FROM chats c
LEFT JOIN preferences ON preferences.chat_id = c.id
WHERE COALESCE(weekly_insights_notifications, true) = true
AND c.status = 'reachable';