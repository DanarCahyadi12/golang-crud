package repository

import (
	"go-crud/internal/entity"
	"gorm.io/gorm"
)

type ProductRepositoryInterface interface {
	Save(product *entity.Product) error
	FindOneById(product *entity.Product, id string) error
	FindMany(products []*entity.Product, offset int, limit int) error
}

type ProductRepository struct {
	Database *gorm.DB
}

func NewProductRepository(database *gorm.DB) *ProductRepository {
	return &ProductRepository{
		Database: database,
	}
}

func (r *ProductRepository) Save(product *entity.Product) error {
	err := r.Database.Create(product).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *ProductRepository) FindOneById(product *entity.Product, id string) error {
	err := r.Database.First(product, "id = ?", id).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *ProductRepository) FindMany(products []*entity.Product, offset int, limit int) error {
	err := r.Database.Limit(limit).Offset(offset).Find(products).Error
	if err != nil {
		return err
	}
	return nil
}
