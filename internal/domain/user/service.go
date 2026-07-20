package user

import (
	"errors"
	"fmt"
	"gotickets/internal/auth"
	"gotickets/internal/domain/user/dto"
)

var ErrInvalidCredentials = errors.New("invalid email or password")

type service struct {
	repo Repository
	jwt  auth.JWTService
}

func NewService(repo Repository, jwt auth.JWTService) *service {
	return &service{
		repo: repo,
		jwt:  jwt,
	}
}

func (s *service) CreateUser(req dto.CreateRequest) (*dto.Response, error) {

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
		return nil, fmt.Errorf("Failed to create user: %w", err)
	}

	// generate token
	accessToken, refreshToken, err := s.jwt.GenerateToken(user.ID, user.Name, user.Email)
	if err != nil {
		return nil, fmt.Errorf("Failed to generate token: %w", err)
	}

	user.RefreshToken = refreshToken
	if err := s.repo.UpdateUser(&user); err != nil {
		return nil, fmt.Errorf("Failed to save refresh token: %w", err)
	}

	response := dto.Response{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		CreatedAt:    user.CreatedAt,
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

	accessToken, refreshToken, err := s.jwt.GenerateToken(user.ID, user.Name, user.Email)
	if err != nil {
		return nil, fmt.Errorf("Failed to generate token: %w", err)
	}

	user.RefreshToken = refreshToken
	if err := s.repo.UpdateUser(user); err != nil {
		return nil, fmt.Errorf("Failed to save refresh token: %w", err)
	}

	response := dto.Response{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		CreatedAt:    user.CreatedAt,
	}

	return &response, nil

}

func (s *service) RefreshToken(token string) (*dto.Response, error) {
	_, err := s.jwt.ValidateToken(token)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	user, err := s.repo.GetUserByRefreshToken(token)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid refresh token")
	}

	accessToken, refreshToken, err := s.jwt.GenerateToken(user.ID, user.Name, user.Email)
	if err != nil {
		return nil, fmt.Errorf("Failed to generate token: %w", err)
	}

	user.RefreshToken = refreshToken
	if err := s.repo.UpdateUser(user); err != nil {
		return nil, fmt.Errorf("Failed to save refresh token: %w", err)
	}

	response := dto.Response{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		CreatedAt:    user.CreatedAt,
	}

	return &response, nil
}
