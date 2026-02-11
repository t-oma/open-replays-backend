-- name: ListVideos :many
SELECT * FROM videos
ORDER BY uploaded_at DESC
LIMIT ? OFFSET ?;

-- name: CountVideos :one
SELECT COUNT(*) FROM videos;

-- name: GetVideo :one
SELECT * FROM videos 
WHERE id = ?;

-- name: CreateVideo :one
INSERT INTO videos (id, title, description, extension, duration, uploaded_at)
VALUES (?, ?, ?, ?, ?, ?) 
RETURNING *;

-- name: UpdateVideoMetadata :exec
UPDATE videos
SET duration = ?
WHERE id = ?;

-- name: DeleteVideo :exec
DELETE FROM videos
WHERE id = ?;
