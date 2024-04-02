package repository

import (
	"go-crud/internal/entity"
	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	Save(user *entity.User) error
	FindOneByEmail(email string) (*entity.User, error)
	FindOneById(user *entity.User, id string) error
	DeleteOneById(id string) error
}
type UserRepository struct {
	Database *gorm.DB
}

func NewUserRepository(database *gorm.DB) *UserRepository {
	return &UserRepository{
		Database: database,
	}
}

func (r *UserRepository) Save(user *entity.User) error {
	err := r.Database.Create(user).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) FindOneByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.Database.First(&user, "email = ?", email).Error
	if err != nil {

		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindOneById(user *entity.User, id string) error {
	err := r.Database.First(user, id).Error
	if err != nil {
		return err
	}

	return nil

}

func (r *UserRepository) DeleteOneById(id string) error {
	err := r.Database.Delete(&entity.User{}, id).Error
	if err != nil {
		return err
	}

	return nil
}
