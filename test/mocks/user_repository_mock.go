package mocks

import (
	"github.com/stretchr/testify/mock"
	"go-crud/internal/entity"
	"go-crud/internal/models"
)

type UserRepositoryMock struct {
	Mock mock.Mock
}

func NewRepositoryMock() *UserRepositoryMock {
	return &UserRepositoryMock{
		Mock: mock.Mock{},
	}
}

func (r *UserRepositoryMock) Save(user *entity.User) error {
	args := r.Mock.Called(user)
	if err := args.Error(0); err != nil {
		return err.(*models.ErrorResponse)
	}

	return nil

}

func (r *UserRepositoryMock) FindOneByEmail(email string) (*entity.User, error) {
	args := r.Mock.Called(email)
	err := args.Error(1)
	if err != nil {
		return nil, args.Error(1).(*models.ErrorResponse)
	}
	return args.Get(0).(*entity.User), nil

}

func (r *UserRepositoryMock) FindOneById(user *entity.User, id string) error {
	args := r.Mock.Called()
	err := args.Error(0)
	if err != nil {
		return args.Error(0)
	}
	return nil
}

func (r *UserRepositoryMock) DeleteOneById(id string) error {
	args := r.Mock.Called()
	err := args.Error(0)
	if err != nil {
		return args.Error(0)
	}
	return nil
}
