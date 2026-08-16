-- name: InsertThread :one
INSERT INTO threads (id, guild_id)
VALUES (:id, :guild_id)
RETURNING *;

-- name: InsertMessage :one
INSERT INTO messages (id, thread_id)
VALUES (:id, :thread_id)
RETURNING *;

-- name: ListMessagesByThread :many
SELECT *
FROM messages
WHERE thread_id = :thread_id
ORDER BY id ASC;
