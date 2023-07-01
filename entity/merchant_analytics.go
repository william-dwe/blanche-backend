package entity

import (
	"time"

	"gorm.io/gorm"
)

type MerchantDailyAnalyticsHist struct {
	ID            uint `gorm:"primary_key"`
	DatePartition time.Time
	Domain        string

	Revenue     float64
	TrxCount    int
	AvgReview   float64
	CountReview int
	OAD         float64
	OSD         float64

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
