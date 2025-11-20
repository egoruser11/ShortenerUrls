package storage

import "errors"

var (
	ErrorURLNotFound = errors.New("url not found")
	ErrURLExists     = errors.New("url already exists")
)
