package model

import (
	"database/sql"
	"gorm.io/gorm"
)

type Link struct {
	gorm.Model
	ID        int64
	Name      string       `gorm:"type:text"`
	Url       string       `gorm:"type:text"`
	CreatedAt sql.NullTime `gorm:"type:TIMESTAMP"`
	ParsedAt  sql.NullTime `gorm:"type:TIMESTAMP NULL"`
	Items     []*Item      `gorm:"foreignKey:LinkID"`
}

func (Link) TableName() string {
	return "links"
}
