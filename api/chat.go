package api

import (
	"github.com/MarkSmersh/nil-chat/db/repo"
	"github.com/gofiber/fiber/v3"
)

func (s Server) ChatGetPreviews(c fiber.Ctx) error {
	userId := c.Locals("userId").(int)

	msgs, err := s.DB.WithAuth(userId).Chat().GetPreviews()

	if err != nil {
		return err
	}

	return c.JSON(msgs)
}

func (s Server) ChatFromUser(c fiber.Ctx) error {
	userId := c.Locals("userId").(int)

	var req repo.ChatFromUser

	err := c.Bind().JSON(&req)

	msgs, err := s.DB.WithAuth(userId).Chat().FromUser(req)

	if err != nil {
		return err
	}

	return c.JSON(msgs)
}

func (s Server) ChatGetAll(c fiber.Ctx) error {
	userId := c.Locals("userId").(int)

	msgs, err := s.DB.WithAuth(userId).Chat().GetAll()

	if err != nil {
		return err
	}

	return c.JSON(msgs)
}

func (s Server) ChatGetHistory(c fiber.Ctx) error {
	userId := c.Locals("userId").(int)

	var req repo.ChatGetHistory

	err := c.Bind().JSON(&req)

	msgs, err := s.DB.WithAuth(userId).Chat().GetHistory(req)

	if err != nil {
		return err
	}

	return c.JSON(msgs)
}

func (s Server) ChatJoinGlobal(c fiber.Ctx) error {
	userId := c.Locals("userId").(int)

	msgs, err := s.DB.WithAuth(userId).Chat().JoinGlobal()

	if err != nil {
		return err
	}

	return c.JSON(msgs)
}
