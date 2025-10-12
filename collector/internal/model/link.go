package model

import (
	"errors"
	"time"
)

var (
	ErrLinkNotFound = errors.New("link not found")
)

type Link struct {
	ID         int64
	Name       string
	Url        string
	CreatedAt  time.Time
	ParsedAt   *time.Time
	Items      []*Item
	ItemsCount int
}
