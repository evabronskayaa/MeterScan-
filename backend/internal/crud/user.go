package crud

import (
	"backend/internal/errors"
	"backend/internal/schema"
	"backend/internal/util"
	errors2 "errors"
	"gorm.io/gorm"
)

const (
	ErrNotFoundUser   errors.SimpleError = "Пользователь не найден"
	ErrDuplicateEmail errors.SimpleError = "Данная почта уже используется"
	ErrSaveUser       errors.SimpleError = "Произошла ошибка при сохранении пользователя"
)

func FindUser(db *gorm.DB, query string, args ...interface{}) (*schema.User, error) {
	var user *schema.User
	if db.Where(query, args...).First(&user).Error != nil {
		return nil, ErrNotFoundUser
	}
	return user, nil
}

func CreateUser(db *gorm.DB, email, password string) (*schema.User, error) {
	_, err := FindUser(db, "email = ?", email)
	if err == nil {
		return nil, ErrDuplicateEmail
	} else if !errors2.Is(err, ErrNotFoundUser) {
		return nil, err
	}
	hashPassword, err := util.HashPassword(password)
	if err != nil {
		return nil, err
	}
	user := &schema.User{
		Email:    email,
		Password: hashPassword,
	}
	if db.Create(&user).Error != nil {
		return nil, ErrSaveUser
	}
	return user, nil
}
