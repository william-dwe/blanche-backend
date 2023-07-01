package entity

import (
	"time"

	"gorm.io/gorm"
)

type MarketplaceDailyAnalyticsHist struct {
	ID                 uint `gorm:"primary_key"`
	DatePartition      time.Time
	CountDailyMau      int
	CountDailyMtu      int
	UserConversionRate float64
	Revenue            float64
	TrxCount           int
	AvgReview          float64
	CountReview        int

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
