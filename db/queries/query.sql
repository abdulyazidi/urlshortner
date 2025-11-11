-- name: GetOriginalURL :one
SELECT original_url FROM urls
WHERE id = ? LIMIT 1;

-- name: CreateURL :one
INSERT INTO urls (
  id, original_url
) VALUES (
  ?, ?
)
RETURNING *;