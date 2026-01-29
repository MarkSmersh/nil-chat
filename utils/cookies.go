package utils

import "github.com/gofiber/fiber/v3"

func AssignTokenToCookies(c fiber.Ctx, accessToken string) {
	c.Cookie(&fiber.Cookie{
		Name:     "access-token",
		Value:    accessToken,
		Path:     "/",
		Secure:   true,
		HTTPOnly: true,
	})
}
