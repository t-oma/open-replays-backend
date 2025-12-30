-- name: ListVideos :many
SELECT * FROM videos;

-- name: GetVideo :one
SELECT * FROM videos 
WHERE filename = ?;

-- name: CreateVideo :one
INSERT INTO videos (title, description, filename, extension, uploaded_at) 
VALUES (?, ?, ?, ?, ?) 
RETURNING *;

-- name: DeleteVideo :exec
DELETE FROM videos 
WHERE title = ?;
