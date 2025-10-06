package model

import (
	"database/sql"
	"errors"
	"time"
)

var (
	ErrLinkNotFound = errors.New("link not found")
)

type Link struct {
	ID        int64
	Name      string
	Url       string
	MinPrice  sql.NullInt64
	MaxPrice  sql.NullInt64
	CreatedAt time.Time
	ParsedAt  *time.Time
	Items     []*Item
}
