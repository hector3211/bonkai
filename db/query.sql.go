// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: query.sql

package db

import (
	"context"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (id,name,email) VALUES (?,?,?) RETURNING id, name, email
`

type CreateUserParams struct {
	ID    interface{}
	Name  string
	Email string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.ID, arg.Name, arg.Email)
	var i User
	err := row.Scan(&i.ID, &i.Name, &i.Email)
	return i, err
}

const deleteuser = `-- name: Deleteuser :exec
DELETE FROM users
WHERE id = ?
`

func (q *Queries) Deleteuser(ctx context.Context, id interface{}) error {
	_, err := q.db.ExecContext(ctx, deleteuser, id)
	return err
}

const getUser = `-- name: GetUser :one
SELECT id, name, email FROM users
WHERE id = ? LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, id interface{}) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, id)
	var i User
	err := row.Scan(&i.ID, &i.Name, &i.Email)
	return i, err
}

const listUsers = `-- name: ListUsers :many
SELECT id, name, email FROM users
`

func (q *Queries) ListUsers(ctx context.Context) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, listUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(&i.ID, &i.Name, &i.Email); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateuser = `-- name: Updateuser :exec
UPDATE users
set name = ?,
email = ?
WHERE id = ?
`

type UpdateuserParams struct {
	Name  string
	Email string
	ID    interface{}
}

func (q *Queries) Updateuser(ctx context.Context, arg UpdateuserParams) error {
	_, err := q.db.ExecContext(ctx, updateuser, arg.Name, arg.Email, arg.ID)
	return err
}
