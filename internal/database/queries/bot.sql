-- name: ListTrackedChannels :many
SELECT id
FROM tracked_channels
WHERE id IN (sqlc.slice('ids'));

-- name: SetChannelTracked :exec
INSERT INTO tracked_channels (id, guild_id)
VALUES (?, ?)
ON CONFLICT (id) DO NOTHING;

-- name: SetChannelUntracked :exec
DELETE FROM tracked_channels
WHERE id = ?;
