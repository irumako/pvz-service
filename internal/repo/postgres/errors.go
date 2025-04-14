package postgres

import "errors"

var (
	ErrDuplicateKey   = errors.New("duplicate key")
	ErrRecordNotFound = errors.New("record not found")
	ErrTransaction    = errors.New("transaction failed")
)
