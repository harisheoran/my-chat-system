package postgre

import (
	"github.com/harisheoran/my-chat-system/pkg/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserController struct {
	DbConnection *gorm.DB
}

// create a new user in the database
func (uc *UserController) CreateNewUser(newUser model.User) error {

	// hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 12)
	if err != nil {
		return err
	}
	newUser.Password = string(hashedPassword)

	// insert user
	result := uc.DbConnection.Create(&newUser)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

// check user exists or not by checking email exists or not in the database
func (uc *UserController) CheckUserExist(email string) (bool, error) {
	var user model.User

	// check for existing user with email
	query := uc.DbConnection.Where("email = ?", email).First(&user)

	if query.Error == gorm.ErrRecordNotFound {
		return false, nil
	}

	if query.Error != nil {
		return false, query.Error
	}

	return true, nil
}

// check the credability of password of the user
func (uc *UserController) Authenticate(email, password string) (model.User, error) {
	user := model.User{}
	result := uc.DbConnection.Where("Email= ?", email).First(&user)

	if result.Error != nil {
		return user, result.Error
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err == bcrypt.ErrMismatchedHashAndPassword {
		return user, gorm.ErrInvalidData
	} else if err != nil {
		return user, err
	}

	return user, nil
}
