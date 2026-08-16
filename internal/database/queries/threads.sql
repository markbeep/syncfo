-- name: InsertThread :one
INSERT INTO threads (id, guild_id)
VALUES (:id, :guild_id)
RETURNING *;

-- name: ListUnbuiltThreads :many
SELECT DISTINCT thread_id
FROM messages
WHERE NOT built
LIMIT ?
OFFSET ?;

-- name: ListUnbuiltMessages :many
SELECT *
FROM messages
WHERE NOT built AND thread_id = ?
LIMIT ?
OFFSET ?;

-- name: InsertMessage :one
INSERT INTO messages (id, thread_id)
VALUES (:id, :thread_id)
RETURNING *;

-- name: SetMessageBuilt :exec
UPDATE messages
SET built = TRUE
WHERE id = ?;

-- name: SetMessageLastEdited :exec
UPDATE messages
SET last_edited = ?
WHERE id = ?;

-- name: ListMessagesByThread :many
SELECT *
FROM messages
WHERE thread_id = ?
ORDER BY id ASC
LIMIT ?
OFFSET ?;
