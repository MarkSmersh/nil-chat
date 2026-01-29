package api

import (
	"errors"
	"time"

	"github.com/MarkSmersh/nil-chat/db/repo"
	"github.com/MarkSmersh/nil-chat/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

func (s Server) AuthRegisterPost(c fiber.Ctx) error {
	req := new(repo.UserRegister)

	if err := c.Bind().JSON(&req); err != nil {
		return fiber.NewError(400, "Invalid request")
	}

	if req.Password == "" || req.Username == "" || req.DisplayName == "" {
		return fiber.NewError(400, "Required fields: username, password, displayName")
	}

	accessToken, err := s.DB.Auth().Register(*req)

	if err != nil {
		return err
	}

	utils.AssignTokenToCookies(c, accessToken)

	return c.Status(201).SendString("Successfully registered")
}

func (s Server) AuthLoginPost(c fiber.Ctx) error {
	req := new(repo.UserLogin)

	if err := c.Bind().JSON(&req); err != nil {
		return fiber.NewError(400, "Invalid request")
	}

	accessToken, err := s.DB.Auth().Login(*req)

	if err != nil {
		return err
	}

	utils.AssignTokenToCookies(c, accessToken)

	return c.SendString("Successfully logged in")
}

func (s Server) AuthLogoutGet(c fiber.Ctx) error {
	accessToken := c.Cookies("access-token")

	if accessToken == "" {
		return fiber.ErrUnauthorized
	}

	iss, sub, err := utils.DecodeAccessToken(accessToken)

	if err != nil && !errors.Is(err, jwt.ErrTokenExpired) {
		return fiber.ErrUnauthorized
	}

	err = s.DB.Auth().RevokeRefreshToken(sub, iss)

	if err != nil {
		return fiber.ErrInternalServerError
	}

	c.Cookie(&fiber.Cookie{
		Name:     "access-token",
		Value:    "",
		Expires:  time.Now().Add(-69 * time.Second),
		Path:     "/",
		Secure:   true,
		HTTPOnly: true,
	})

	return c.SendString("Successfully logged out")
}
