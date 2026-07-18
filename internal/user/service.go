package user

import "gotickets/internal/user/dto"

type service struct {
	repo Repository
}

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

	err = s.repo.CreateUser(&user)

	if err != nil {
		return &dto.Response{}, err
	}

	response := dto.Response{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
	}

	return &response, nil

}
