package user

import (
	"errors"
	"gotickets/internal/user/dto"
)

type service struct {
	repo Repository
}

var ErrInvalidCredentials = errors.New("invalid email or password")

func NewService(repo Repository) *service {
	return &service{
		repo: repo,
	}
}

func (s service) CreateUser(req dto.CreateRquest) (*dto.Response, error) {

	user := User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	err := user.HashPassword(user.Password)
	if err != nil {
		return &dto.Response{}, err
	}

	// generate token

	err = s.repo.CreateUser(&user)

	if err != nil {
		return &dto.Response{}, err
	}

	response := dto.Response{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	return &response, nil

}

func (s *service) LoginUser(req dto.LoginRequest) (*dto.Response, error) {

	user, err := s.repo.GetUserByEmail(&User{
		Email: req.Email,
	})

	if err != nil {
		return &dto.Response{}, err
	}

	if user == nil {
		return nil, ErrInvalidCredentials
	}

	if err := user.CheckPassword(req.Password); err != nil {
		return nil, ErrInvalidCredentials
	}

	err = s.repo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	response := dto.Response{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	return &response, nil

}
