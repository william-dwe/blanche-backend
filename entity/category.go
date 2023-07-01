package entity

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID            uint `gorm:"primaryKey"`
	Name          string
	Slug          string `gorm:"uniqueIndex"`
	ImageUrl      string
	ParentId      uint `gorm:"default:null"`
	Parent        *Category
	GrandparentId uint `gorm:"default:null"`
	Grandparent   *Category
	Children      []Category `gorm:"foreignKey:ParentId"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
