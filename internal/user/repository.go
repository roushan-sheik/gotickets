package user

import (
	"errors"

	"gorm.io/gorm"
)

var ErrorAlreadyExist = errors.New("user with this email aready exists")

type Repository interface {
	CreateUser(user *User) error
	GetUserByEmail(user *User) (*User, error)
	UpdateUser(user *User) error
	GetUserByRefreshToken(token string) (*User, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r repository) CreateUser(user *User) error {
	result := r.db.Create(user)

	if result.Error != nil {

		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return ErrorAlreadyExist
		}

		return result.Error
	}
	return nil

}

func (r repository) GetUserByEmail(user *User) (*User, error) {
	result := r.db.Where("email = ?", user.Email).First(user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return user, nil
}

func (r repository) UpdateUser(user *User) error {
	result := r.db.Save(user)
	return result.Error
}

func (r repository) GetUserByRefreshToken(token string) (*User, error) {
	var user User
	result := r.db.Where("refresh_token = ?", token).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}
