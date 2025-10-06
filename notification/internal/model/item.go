package model

import (
	"errors"
	"time"
)

var (
	ErrItemNotFound = errors.New("item not found")
)

type Item struct {
	ID          int64
	LinkId      int64
	Uid         string
	Title       string
	Description string
	Url         string
	PreviewUrl  string
	ImageUrls   []string
	Price       int
	NeedNotify  bool
	CreatedAt   time.Time
}

type ItemEvent struct {
	Title       string
	Description string
	Url         string
	PreviewUrl  string
	Price       int
	CreatedAt   time.Time
}
