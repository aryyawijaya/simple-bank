package entity

import "errors"

var (
	ErrUnique          = errors.New("any field in your resource already existed")
	ErrNotFound        = errors.New("your requested resource is not found")
	ErrTokenExpired    = errors.New("token has expired")
	ErrTokenInvalid    = errors.New("token is invalid")
	ErrSessionBlocked  = errors.New("blocked session")
	ErrSessionInvalid  = errors.New("invalid session")
	ErrSessionExpired  = errors.New("expired session")
	ErrPasswordInvalid = errors.New("invalid password")
)
