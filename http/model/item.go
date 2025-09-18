package model

import (
	"errors"
	"gorm.io/gorm"
)

var (
	ErrItemNotFound = errors.New("item not found")
)

type Item struct {
	gorm.Model
	LinkID      uint
	Uid         string `gorm:"type:varchar(255)"`
	Title       string `gorm:"type:text"`
	Description string `gorm:"type:text"`
	Url         string `gorm:"type:text"`
	PreviewUrl  string `gorm:"type:text"`
	Price       int    `gorm:"type:int"`
	NeedNotify  bool   `gorm:"type:boolean"`
}

func (Item) TableName() string {
	return "items"
}
