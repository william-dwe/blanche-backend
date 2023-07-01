package entity

import (
	"time"

	"gorm.io/gorm"
)

type SlpAccount struct {
	ID         uint `gorm:"primary_key"`
	UserID     uint
	CardNumber string
	ActiveDate time.Time
	NameOnCard string
	IsDefault  bool

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
