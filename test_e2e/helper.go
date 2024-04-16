package test_e2e

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"go-crud/internal/entity"
	"go-crud/internal/models"
	"gorm.io/gorm"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
)

var email = "danar@gmail.com"

func SignIn(request models.SignInRequest) (*models.Response[*models.AuthResponse], error) {
	bodyJson, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req := httptest.NewRequest(http.MethodPost, "/auth/signin", strings.NewReader(string(bodyJson)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	response, err := App.Fiber.Test(req)
	if err != nil {
		return nil, err
	}

	var authResponse models.Response[*models.AuthResponse]

	bodyByte, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bodyByte, &authResponse)
	if err != nil {
		return nil, err
	}

	return &authResponse, nil

}
func DeleteUserIfExits() {
	userFound, err := UserRepository.FindOneByEmail(email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		panic(err)
	}

	if userFound != nil {
		err := UserRepository.DeleteOneById(userFound.Id)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			panic(err)
		}
	}

}

var user = &entity.User{
	Name:     "Danar",
	Email:    "danar@gmail.com",
	Password: "$2a$10$aOySpFRuA2uE8gGNNCuAleiBvNRyMJpZuyhZ21kf/Tpy5c8KHNRTe",
}

func CreateUser(user *entity.User) {
	user.Id = uuid.New().String()
	userFound, err := UserRepository.FindOneByEmail(user.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		panic(err)
	}

	if userFound != nil {
		err := UserRepository.DeleteOneById(userFound.Id)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			panic(err)
		}
	}

	err = UserRepository.Save(user)
	if err != nil {
		panic(err)
	}
}

func DeleteUser(id string) {
	err := UserRepository.DeleteOneById(id)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		panic(err)
	}
}

func CreateManyProduct(userID string) error {
	var products []entity.Product
	for i := 1; i <= 50; i++ {
		products = append(products, entity.Product{
			Id:     uuid.New().String(),
			Name:   fmt.Sprintf("Product %d", i),
			Price:  20000,
			Stock:  120,
			UserId: userID,
		})
	}
	err := Database.Create(&products).Error
	if err != nil {
		return err
	}
	return nil
}
