package api

import "github.com/gofiber/fiber/v3"

func (s *Server) HealthyGet(c fiber.Ctx) error {
	return c.Status(200).SendString("Server is healthy")
}
