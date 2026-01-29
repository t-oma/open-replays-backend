-- name: ListVideos :many
SELECT * FROM videos;

-- name: GetVideo :one
SELECT * FROM videos 
WHERE id = ?;

-- name: CreateVideo :one
INSERT INTO videos (id, title, description, filename, extension, duration, thumbnail, uploaded_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?) 
RETURNING *;

-- name: DeleteVideo :exec
DELETE FROM videos
WHERE id = ?;
