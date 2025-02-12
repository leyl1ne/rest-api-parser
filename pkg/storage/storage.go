package storage

import "errors"

var (
	ErrSongNotFound   = errors.New("song not found")
	ErrSongExists     = errors.New("song exists")
	ErrArtistNotFound = errors.New("artist not found")
	ErrArtistExists   = errors.New("artist exists")
)
