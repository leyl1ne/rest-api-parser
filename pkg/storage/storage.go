package storage

import "errors"

var (
	ErrSongNotFound = errors.New("song not found")
	ErrSongExists   = errors.New("url exists")
)
