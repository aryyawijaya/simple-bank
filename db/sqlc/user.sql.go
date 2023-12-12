// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: user.sql

package mydb

import (
	"context"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
  "username",
  "hashedPassword",
  "fullName",
  "email"
) VALUES (
  $1, $2, $3, $4
)
RETURNING username, "hashedPassword", "fullName", email, "passwordChangedAt", "createdAt"
`

type CreateUserParams struct {
	Username       string `json:"username"`
	HashedPassword string `json:"hashedPassword"`
	FullName       string `json:"fullName"`
	Email          string `json:"email"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.Username,
		arg.HashedPassword,
		arg.FullName,
		arg.Email,
	)
	var i User
	err := row.Scan(
		&i.Username,
		&i.HashedPassword,
		&i.FullName,
		&i.Email,
		&i.PasswordChangedAt,
		&i.CreatedAt,
	)
	return i, err
}

const getUser = `-- name: GetUser :one
SELECT username, "hashedPassword", "fullName", email, "passwordChangedAt", "createdAt" FROM users
WHERE "username" = $1
LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, username)
	var i User
	err := row.Scan(
		&i.Username,
		&i.HashedPassword,
		&i.FullName,
		&i.Email,
		&i.PasswordChangedAt,
		&i.CreatedAt,
	)
	return i, err
}
