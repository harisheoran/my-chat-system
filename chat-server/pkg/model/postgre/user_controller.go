package postgre

import (
	"strings"

	"github.com/harisheoran/my-chat-system/pkg/model"
	"gorm.io/gorm"
)

type UserController struct {
	DbConnection *gorm.DB
}

func (uc *UserController) CheckUserExists(email string) (interface{}, error) {
	var user model.User

	// check for existing user with email
	query := uc.DbConnection.Where("email = ?", email).First(&user)

	if query.Error != nil {
		return nil, query.Error
	}

	// send existing user
	response := map[string]interface{}{
		"username": user.Username,
		"email":    user.Email,
		"password": user.Password,
		"type":     "login",
	}
	return response, nil

}

func (uc *UserController) CreateNewUser(name string, email string, password string) (interface{}, error) {
	newUser := model.User{
		Name:     name,
		Username: strings.Split(email, "@")[0],
		Email:    email,
		Password: password,
	}

	query := uc.DbConnection.Create(&newUser)

	if query.Error != nil {
		return nil, query.Error
	}

	// send user data
	response := map[string]interface{}{
		"username": newUser.Username,
		"email":    newUser.Email,
		"name":     newUser.Name,
		"type":     "signup",
	}
	return response, nil
}
