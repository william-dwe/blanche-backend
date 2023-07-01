package entity

type PromotionType struct {
	ID   uint `gorm:"primaryKey"`
	Name string
}
