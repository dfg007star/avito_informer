package model

import (
	"database/sql"
	"errors"
	"gorm.io/gorm"
)

var (
	ErrLinkNotFound = errors.New("link not found")
)

type Link struct {
	gorm.Model
	Name      string       `gorm:"type:text"`
	Url       string       `gorm:"type:text"`
	CreatedAt sql.NullTime `gorm:"type:TIMESTAMP"`
	ParsedAt  sql.NullTime `gorm:"type:TIMESTAMP NULL"`
}

func (Link) TableName() string {
	return "links"
}
