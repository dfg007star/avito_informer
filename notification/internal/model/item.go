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
	IsNotify    bool
	CreatedAt   time.Time

	CategoryTitle string
	LinkUrl       string
}

type ItemEvent struct {
	Title         string
	CategoryTitle string
	Description   string
	Url           string
	PreviewUrl    string
	Price         int
	CreatedAt     time.Time
	LinkUrl       string
}
