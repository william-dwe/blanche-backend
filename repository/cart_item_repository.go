package repository

import (
	"time"

	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/entity"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type CartItemRepository interface {
	AddCartItem(cartItem *entity.CartItem) error
	GetCartItems(userId uint) ([]entity.CartItem, error)
	GetSelectedCartItems(userId uint) ([]entity.CartItem, error)
	GetCartItem(userId uint, cartItemId uint) (*entity.CartItem, error)
	GetSameExistingCartItem(userId uint, cartItem *entity.CartItem) (*entity.CartItem, error)
	UpdateCartItem(cartItem *entity.CartItem) error
	UpdateCartItems(userId uint, cartItemIds []uint, newCartItem *map[string]interface{}) error
	DeleteCartItemByCartId(userId uint, cartId uint) error
	DeleteSelectedCartItem(userId uint) error
	DeleteCartItemsByProductVariantIds(userId uint, productVariantIds []uint) error
}
type CartItemRepositoryConfig struct {
	DB *gorm.DB
}

type cartItemRepositoryImpl struct {
	db *gorm.DB
}

func NewCartItemRepository(c CartItemRepositoryConfig) CartItemRepository {
	return &cartItemRepositoryImpl{
		db: c.DB,
	}
}

func (r *cartItemRepositoryImpl) AddCartItem(cartItem *entity.CartItem) error {
	err := r.db.Create(&cartItem).Error

	if err != nil {
		pgErr := err.(*pgconn.PgError)

		if pgErr.Code == "23503" {
			if pgErr.ConstraintName == "cart_items_product_id_fkey" {
				return domain.ErrAddCartProductNotExist
			}
			if pgErr.ConstraintName == "cart_items_user_id_fkey" {
				return domain.ErrAddCartUserNotExist
			}
			if pgErr.ConstraintName == "cart_items_product_id_fkey" {
				return domain.ErrAddCartVariantItemNotExist
			}
		}

		log.Error().Msgf("Error create cart item: %v", pgErr.Message)
		return domain.ErrRegister
	}

	return err
}

func (r *cartItemRepositoryImpl) GetCartItems(userId uint) ([]entity.CartItem, error) {
	var cartItems []entity.CartItem
	err := r.db.
		Preload("Product", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped()
		}).
		Preload("Product.ProductPromotion.Promotion", func(db *gorm.DB) *gorm.DB {
			return db.Where("start_at <= ?", time.Now()).Where("end_at >= ?", time.Now())
		}).
		Preload("Product.Merchant").
		Preload("Product.ProductImages").
		Preload("VariantItem", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped()
		}).
		Preload("VariantItem.VariantSpec", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped()
		}).
		Where("user_id = ?", userId).
		Order("created_at desc").
		Find(&cartItems).
		Error
	if err != nil {
		return nil, domain.ErrGetCartInternalError
	}

	return cartItems, nil
}

func (r *cartItemRepositoryImpl) GetSelectedCartItems(userId uint) ([]entity.CartItem, error) {
	const is_checked = true
	var cartItems []entity.CartItem
	err := r.db.
		Preload("Product", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped()
		}).
		Preload("Product.Merchant").
		Preload("Product.ProductImages").
		Preload("VariantItem", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped()
		}).
		Where("user_id = ? AND is_checked = ?", userId, is_checked).
		Order("created_at desc").
		Find(&cartItems).
		Error
	if err != nil {
		return nil, domain.ErrGetCartInternalError
	}

	return cartItems, nil
}

func (r *cartItemRepositoryImpl) GetCartItem(userId uint, cartItemId uint) (*entity.CartItem, error) {
	var cartItem entity.CartItem
	err := r.db.
		Preload("Product", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped()
		}).
		Preload("Product.Merchant").
		Preload("Product.ProductImages").
		Preload("VariantItem", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped()
		}).
		Where("user_id = ? AND id = ?", userId, cartItemId).
		Order("created_at desc").
		First(&cartItem).
		Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrInvalidCartItemID
		}
		return nil, domain.ErrGetCartInternalError
	}
	return &cartItem, nil
}

func (r *cartItemRepositoryImpl) GetSameExistingCartItem(userId uint, cartItem *entity.CartItem) (*entity.CartItem, error) {
	var similarCartItems entity.CartItem
	err := r.db.
		Preload("Product").
		Preload("Product.Merchant").
		Preload("Product.ProductImages").
		Preload("VariantItem").
		Where("user_id = ? AND product_id = ? AND variant_item_id = ?", userId, cartItem.ProductId, cartItem.VariantItemId).
		Order("created_at desc").
		First(&similarCartItems).
		Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, domain.ErrGetCartInternalError
	}
	return &similarCartItems, nil
}

func (r *cartItemRepositoryImpl) UpdateCartItem(cartItem *entity.CartItem) error {
	q := r.db.
		Model(&entity.CartItem{}).
		Where("id = ? AND user_id = ?",
			cartItem.ID,
			cartItem.UserId,
		).
		Updates(&cartItem)
	err := q.Error
	if err != nil {
		pgErr := err.(*pgconn.PgError)

		if pgErr.Code == "23503" {
			if pgErr.ConstraintName == "cart_items_product_id_fkey" {
				return domain.ErrAddCartProductNotExist
			}
			if pgErr.ConstraintName == "cart_items_user_id_fkey" {
				return domain.ErrAddCartUserNotExist
			}
			if pgErr.ConstraintName == "cart_items_product_id_fkey" {
				return domain.ErrAddCartVariantItemNotExist
			}
		}

		log.Error().Msgf("Error update cart item: %v", pgErr.Message)
		return domain.ErrRegister
	}
	if q.RowsAffected == 0 {
		return domain.ErrInvalidCartItemID
	}

	return nil
}

func (r *cartItemRepositoryImpl) UpdateCartItems(userId uint, cartItemIds []uint, newCartItem *map[string]interface{}) error {
	q := r.db.
		Model(&entity.CartItem{}).
		Where("id in (?) AND user_id = ?", cartItemIds, userId).
		Updates(&newCartItem)
	err := q.Error
	if err != nil {
		pgErr := err.(*pgconn.PgError)

		if pgErr.Code == "23503" {
			if pgErr.ConstraintName == "cart_items_product_id_fkey" {
				return domain.ErrAddCartProductNotExist
			}
			if pgErr.ConstraintName == "cart_items_user_id_fkey" {
				return domain.ErrAddCartUserNotExist
			}
			if pgErr.ConstraintName == "cart_items_product_id_fkey" {
				return domain.ErrAddCartVariantItemNotExist
			}
		}

		log.Error().Msgf("Error update cart items: %v", pgErr.Message)
		return domain.ErrRegister
	}
	if q.RowsAffected == 0 {
		return domain.ErrInvalidCartItemID
	}

	return nil
}

func (r *cartItemRepositoryImpl) DeleteSelectedCartItem(userId uint) error {
	const is_checked = true
	q := r.db.
		Where("user_id = ? AND is_checked = ?", userId, is_checked).
		Delete(&entity.CartItem{})
	err := q.Error
	if err != nil {
		return domain.ErrDeleteCartInternalError
	}
	if q.RowsAffected == 0 {
		return domain.ErrInvalidCartItemID
	}

	return nil
}

func (r *cartItemRepositoryImpl) DeleteCartItemByCartId(userId uint, cartId uint) error {
	q := r.db.
		Where("user_id = ? AND id = ?", userId, cartId).
		Delete(&entity.CartItem{})
	err := q.Error
	if err != nil {
		return domain.ErrDeleteCartInternalError
	}
	if q.RowsAffected == 0 {
		return domain.ErrInvalidSelectedCartItemId
	}
	return nil
}

func (r *cartItemRepositoryImpl) DeleteCartItemsByProductVariantIds(userId uint, productVariantIds []uint) error {
	q := r.db.
		Where("user_id = ? AND variant_item_id in (?)", userId, productVariantIds).
		Delete(&entity.CartItem{})
	err := q.Error
	if err != nil {
		return domain.ErrDeleteCartInternalError
	}
	if q.RowsAffected == 0 {
		return domain.ErrInvalidSelectedCartItemId
	}
	return nil
}
