package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	jwtSecretKey    = "your_secret_key"
	tokenExpiryTime = 7 * 24 * time.Hour // after 24 hours token will expire
)

type JwtCustomClaims struct {
	UserID uint   `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type JWTService interface {
	GenerateToken(userId uint, name string, email string) (string, error)
	ValidateToken(tokenString string) (*JwtCustomClaims, error)
}

type jwtService struct {
	jwtSecretKey []byte
	jwtExpiry    time.Duration
}

func NewJWTService(secretKey string) JWTService {

	if secretKey == "" {
		secretKey = jwtSecretKey
	}
	return &jwtService{
		jwtSecretKey: []byte(secretKey),
		jwtExpiry:    tokenExpiryTime,
	}
}

func (js *jwtService) GenerateToken(userId uint, name string, email string) (string, error) {
	claims := &JwtCustomClaims{
		UserID: userId,
		Name:   name,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(js.jwtExpiry)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(js.jwtSecretKey)
}

func (js *jwtService) ValidateToken(tokenString string) (*JwtCustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return js.jwtSecretKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*JwtCustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
