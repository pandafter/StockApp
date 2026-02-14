package shared

import "errors"

var (
	ErrNotFound      = errors.New("record not found")
	ErrUnauthorized  = errors.New("unauthorized")
	ErrForbidden     = errors.New("forbidden")
	ErrBadRequest    = errors.New("bad request")
)
