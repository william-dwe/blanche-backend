package entity

import (
	"time"

	"gorm.io/gorm"
)

type ProductAnalytic struct {
	ID               uint `gorm:"primarykey"`
	Score            float64
	AvgRating        float64
	NumOfReview      int
	NumOfSale        int
	NumOfFavorite    int
	NumOfPendingSale int
	TotalStock       int

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
