package mocks

import (
	"github.com/stretchr/testify/mock"
	"go-crud/internal/entity"
)

type ProductRepositoryMock struct {
	Mock mock.Mock
}

func NewProductRepositoryMock() *ProductRepositoryMock {
	return &ProductRepositoryMock{
		Mock: mock.Mock{},
	}
}

func (r *ProductRepositoryMock) Save(product *entity.Product) error {
	args := r.Mock.Called(product)
	if args.Error(0) != nil {
		return args.Error(0)
	}

	return nil
}
func (r *ProductRepositoryMock) FindOneById(product *entity.Product, id string) error {
	args := r.Mock.Called(product, id)
	err := args.Error(0)
	if err != nil {
		return args.Error(0)
	}

	return nil
}

func (r *ProductRepositoryMock) FindMany(products *[]entity.Product, offset int, limit int) error {
	args := r.Mock.Called(products, offset, limit)
	if args.Error(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (r *ProductRepositoryMock) UpdateById(product entity.Product, productID string) (*entity.Product, error) {
	args := r.Mock.Called(product, productID)
	err := args.Error(1)
	if err != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Product), nil
}

func (r *ProductRepositoryMock) DeleteById(productID string) error {
	args := r.Mock.Called(productID)
	err := args.Error(0)
	if err != nil {
		return err
	}
	return nil
}
func (r *ProductRepositoryMock) Count() (int64, error) {
	args := r.Mock.Called()
	err := args.Error(1)
	if err != nil {
		return -1, err
	}

	return args.Get(0).(int64), nil
}
