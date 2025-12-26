-- name: GetVideosByTitle :many
SELECT * FROM videos 
WHERE title LIKE ?;

-- name: CreateVideo :one
INSERT INTO videos (title, description, filename, uploaded_at) 
VALUES (?, ?, ?, ?) 
RETURNING *;

-- name: DeleteVideo :exec
DELETE FROM videos 
WHERE title = ?;
