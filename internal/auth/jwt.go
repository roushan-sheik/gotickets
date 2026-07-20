package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	jwtAccessSecretKey     = "your_access_secret_key"
	jwtRefreshSecretKey    = "your_refresh_secret_key"
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
	ValidateToken(tokenString string, isRefresh bool) (*JwtCustomClaims, error)
}

type jwtService struct {
	jwtAccessSecretKey  []byte
	jwtRefreshSecretKey []byte
	accessTokenExp      time.Duration
	refreshTokenExp     time.Duration
}

func NewJWTService(accessSecretKey, refreshSecretKey string) JWTService {

	if accessSecretKey == "" {
		accessSecretKey = jwtAccessSecretKey
	}
	if refreshSecretKey == "" {
		refreshSecretKey = jwtRefreshSecretKey
	}
	return &jwtService{
		jwtAccessSecretKey:  []byte(accessSecretKey),
		jwtRefreshSecretKey: []byte(refreshSecretKey),
		accessTokenExp:      accessTokenExpiryTime,
		refreshTokenExp:     refreshTokenExpiryTime,
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
	aToken, err := accessToken.SignedString(js.jwtAccessSecretKey)
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
	rToken, err := refreshToken.SignedString(js.jwtRefreshSecretKey)
	if err != nil {
		return "", "", err
	}

	return aToken, rToken, nil
}

func (js *jwtService) ValidateToken(tokenString string, isRefresh bool) (*JwtCustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		if isRefresh {
			return js.jwtRefreshSecretKey, nil
		}
		return js.jwtAccessSecretKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*JwtCustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
