package repository

import (
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/cache"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/entity"

	"gorm.io/gorm"
)

type ExampleRepository interface {
	Store(entity.Example) error
	CachedStore(input entity.Example) error
	GetByID(uint) (*entity.Example, error)
}

type ExampleRepositoryConfig struct {
	DB  *gorm.DB
	RDB *cache.RDBConnection
}

type exampleRepositoryImpl struct {
	db  *gorm.DB
	rdb *cache.RDBConnection
}

func NewExampleRepository(c ExampleRepositoryConfig) ExampleRepository {
	return &exampleRepositoryImpl{
		db:  c.DB,
		rdb: c.RDB,
	}
}

func (r *exampleRepositoryImpl) Store(input entity.Example) error {
	err := r.db.Create(&input).Error
	if err != nil {
		return domain.ErrCreateExample
	}

	return nil
}

func (r *exampleRepositoryImpl) CachedStore(input entity.Example) error {
	var data entity.Example
	err := r.rdb.GetCache("test", &data)
	if err == nil {
		return nil
	}
	// fetch data from DB here
	err = r.rdb.SetCache("test", entity.Example{Name: "test", Qty: 10}, 10)
	if err != nil {
		return err
	}
	return nil
}

func (r *exampleRepositoryImpl) GetByID(id uint) (*entity.Example, error) {
	var example entity.Example
	err := r.db.First(&example, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrExampleIdNotFound
		}
		//check another error or if is not masked, return internal error
		return nil, domain.ErrGetExample
	}

	return &example, nil
}
