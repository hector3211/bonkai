-- name: GetUser :one
SELECT * FROM users
WHERE id = ? LIMIT 1;


-- name: ListUsers :many 
SELECT * FROM users;


-- name: CreateUser :one
INSERT INTO users (id,name,email) VALUES (?,?,?) RETURNING *;

-- name: Updateuser :exec
UPDATE users
set name = ?,
email = ?
WHERE id = ?;

-- name: Deleteuser :exec
DELETE FROM users
WHERE id = ?;
