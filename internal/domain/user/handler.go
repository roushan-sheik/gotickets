package user

import (
	"gotickets/internal/auth"
	"gotickets/internal/httpresponse"
	"gotickets/internal/domain/user/dto"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
)

type handler struct {
	service *service
}

func NewHandler(service *service) *handler {
	return &handler{
		service: service,
	}
}

// CreateUser godoc
// @Summary      Register a new user
// @Description  Creates a new user account and returns access and refresh tokens.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body      dto.CreateRequest  true  "User Registration Details"
// @Success      201      {object}  dto.Response
// @Failure      400      {object}  httpresponse.Error
// @Router       /api/v1/auth/register [post]
func (h *handler) CreateUser(c *echo.Context) error {
	var req dto.CreateRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Invalid request payload", err.Error()))
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Validation failed", err.Error()))
	}

	response, err := h.service.CreateUser(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Error creating user", err.Error()))
	}

	setTokensInCookies(c, response.AccessToken, response.RefreshToken)

	return c.JSON(http.StatusCreated, response)

}

// LoginUser godoc
// @Summary      Login user
// @Description  Authenticates a user and returns access and refresh tokens.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body      dto.LoginRequest  true  "User Login Details"
// @Success      200      {object}  dto.Response
// @Failure      400      {object}  httpresponse.Error
// @Router       /api/v1/auth/login [post]
func (h *handler) LoginUser(c *echo.Context) error {
	var req dto.LoginRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Invalid request payload", err.Error()))
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Validation failed", err.Error()))
	}

	response, err := h.service.LoginUser(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Error logging in", err.Error()))
	}

	setTokensInCookies(c, response.AccessToken, response.RefreshToken)

	return c.JSON(http.StatusOK, response)
}

// GetMe godoc
// @Summary      Get current user profile
// @Description  Returns the profile of the currently authenticated user based on the JWT token.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200      {object}  auth.JwtCustomClaims
// @Failure      401      {object}  httpresponse.Error
// @Router       /api/v1/users/me [get]
func (h *handler) GetMe(c *echo.Context) error {
	userClaims, ok := c.Get("user").(*auth.JwtCustomClaims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, httpresponse.NewError(http.StatusUnauthorized, "Unauthorized", "User not found in context"))
	}
	return c.JSON(http.StatusOK, userClaims)
}

func setTokensInCookies(c *echo.Context, accessToken, refreshToken string) {
	c.SetCookie(&http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
	})
	c.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Now().Add(30 * 24 * time.Hour),
	})
}

// RefreshToken godoc
// @Summary      Refresh Access Token
// @Description  Generates a new access token using a valid refresh token.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body      object{refresh_token=string}  false  "Refresh Token (Optional if provided in cookies)"
// @Success      200      {object}  dto.Response
// @Failure      401      {object}  httpresponse.Error
// @Router       /api/v1/auth/refresh [post]
func (h *handler) RefreshToken(c *echo.Context) error {
	var token string
	if cookie, err := c.Cookie("refresh_token"); err == nil {
		token = cookie.Value
	} else {
		var req struct {
			RefreshToken string `json:"refresh_token"`
		}
		if bindErr := c.Bind(&req); bindErr == nil {
			token = req.RefreshToken
		}
	}

	if token == "" {
		return c.JSON(http.StatusUnauthorized, httpresponse.NewError(http.StatusUnauthorized, "Missing refresh token", ""))
	}

	response, err := h.service.RefreshToken(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, httpresponse.NewError(http.StatusUnauthorized, "Invalid refresh token", err.Error()))
	}

	setTokensInCookies(c, response.AccessToken, response.RefreshToken)

	return c.JSON(http.StatusOK, response)
}
