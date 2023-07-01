package entity

import (
	"time"

	"gorm.io/gorm"
)

const UserScope = "user"
const UserPaymentEditScope = "user_payment_edit"
const UserCredentialEditScope = "user_credential_edit"

type User struct {
	ID               uint `gorm:"primary_key"`
	RoleId           uint
	Role             Role
	Username         string `gorm:"unique"`
	Email            string `gorm:"unique"`
	Password         string
	FavoriteProducts []Product `gorm:"many2many:user_favorite_products;"`
	UserDetail       UserDetail
	UserAddress      []UserAddress

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
