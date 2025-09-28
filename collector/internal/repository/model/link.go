package model

import (
	"database/sql"
	"gorm.io/gorm"
)

// type Link struct {
//     gorm.Model
//     Name     string       `gorm:"type:text"`
//     Url      string       `gorm:"type:text"`
//     ParsedAt sql.NullTime `gorm:"type:TIMESTAMP NULL"`
//
//     // Relation: one Link has many Items
//     Items []*Item `gorm:"foreignKey:LinkID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
// }

type Link struct {
	gorm.Model
	Name     string       `gorm:"type:text"`
	Url      string       `gorm:"type:text"`
	ParsedAt sql.NullTime `gorm:"type:TIMESTAMP NULL"`
	Items    []*Item      `gorm:"foreignKey:LinkId"`
}

func (Link) TableName() string {
	return "links"
}
