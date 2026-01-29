package api

import (
	"errors"
	"log/slog"

	"github.com/MarkSmersh/nil-chat/db/repo"
	"github.com/MarkSmersh/nil-chat/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

func (s Server) AuthMiddleware(c fiber.Ctx) error {
	accessToken := c.Cookies("access-token")

	accessToken, err := ProcessAccessToken(s.DB.Auth(), accessToken)

	if err != nil {
		return err
	}

	_, sub, _ := utils.DecodeAccessToken(accessToken)

	if c.IsWebSocket() {
		c.Locals("access-token", accessToken)
	}

	c.Locals("userId", sub)

	return c.Next()
}

// checks is the given access token expired, if so - returns a new one
func ProcessAccessToken(repo repo.Auth, accessToken string) (string, error) {
	if accessToken == "" {
		return "", fiber.NewError(401, "The access token is absent")
	}

	iss, sub, err := utils.DecodeAccessToken(accessToken)

	if err == nil {
		return accessToken, nil
	}

	if !errors.Is(err, jwt.ErrTokenExpired) {
		return "", fiber.NewError(400, "The access token is invalid")
	}

	accessToken, err = repo.UpdateAccessToken(sub, iss)

	if err == nil {
		return accessToken, nil
	}

	if errors.Is(err, jwt.ErrTokenExpired) {
		return "", fiber.NewError(401, "Session is expired")
	}

	slog.Error(err.Error())

	return "", fiber.NewError(400, "Unable to update an access token")
}
