package model

import "errors"

var (
	ErrPlayerExists      error = errors.New("player already exists")
	ErrInvalidPlayer     error = errors.New("invalid player")
	ErrRankDoensnotExist error = errors.New("rank doesnt exists")
	ErrNoPlayer          error = errors.New("no player found")
)
