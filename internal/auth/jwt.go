package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	jwtSecretKey           = "your_secret_key"
	accessTokenExpiryTime  = 24 * time.Hour
	refreshTokenExpiryTime = 30 * 24 * time.Hour
)

type JwtCustomClaims struct {
	UserID uint   `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type JWTService interface {
	GenerateToken(userId uint, name string, email string) (string, string, error)
	ValidateToken(tokenString string) (*JwtCustomClaims, error)
}

type jwtService struct {
	jwtSecretKey    []byte
	accessTokenExp  time.Duration
	refreshTokenExp time.Duration
}

func NewJWTService(secretKey string) JWTService {

	if secretKey == "" {
		secretKey = jwtSecretKey
	}
	return &jwtService{
		jwtSecretKey:    []byte(secretKey),
		accessTokenExp:  accessTokenExpiryTime,
		refreshTokenExp: refreshTokenExpiryTime,
	}
}

func (js *jwtService) GenerateToken(userId uint, name string, email string) (string, string, error) {
	// Generate Access Token
	accessClaims := &JwtCustomClaims{
		UserID: userId,
		Name:   name,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(js.accessTokenExp)),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	aToken, err := accessToken.SignedString(js.jwtSecretKey)
	if err != nil {
		return "", "", err
	}

	// Generate Refresh Token
	refreshClaims := &JwtCustomClaims{
		UserID: userId,
		Name:   name,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(js.refreshTokenExp)),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	rToken, err := refreshToken.SignedString(js.jwtSecretKey)
	if err != nil {
		return "", "", err
	}

	return aToken, rToken, nil
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
