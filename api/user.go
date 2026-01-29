package api

import (
	"github.com/MarkSmersh/nil-chat/db/repo"
	"github.com/gofiber/fiber/v3"
)

func (s Server) UserSearch(c fiber.Ctx) error {
	q := c.Params("q")

	users, err := s.DB.User().Search(q)

	if err != nil {
		return err
	}

	return c.JSON(users)
}

func (s Server) UserBlock(c fiber.Ctx) error {
	var req repo.UserBlock

	err := c.Bind().JSON(&req)

	if err != nil {
		return fiber.NewError(400, "Nieprawidlowe żądanie")
	}

	res, err := s.DB.User().Block(req)

	if err != nil {
		return err
	}

	return c.JSON(res)
}

func (s Server) UserUnblock(c fiber.Ctx) error {
	var req repo.UserUnblock

	err := c.Bind().JSON(&req)

	if err != nil {
		return fiber.NewError(400, "Nieprawidlowe żądanie")
	}

	res, err := s.DB.User().Unblock(req)

	if err != nil {
		return err
	}

	return c.JSON(res)
}

// func (s Server) UserRequest(c fiber.Ctx) error {
// 	var req repo.UserRequest
//
// 	err := c.Bind().JSON(&req)
//
// 	if err != nil {
// 		return fiber.NewError(400, "Nieprawidlowe żądanie")
// 	}
//
// 	res, err := s.DB.User().Request(req)
//
// 	if err != nil {
// 		return err
// 	}
//
// 	return c.JSON(res)
// }
//
// func (s Server) UserAccept(c fiber.Ctx) error {
// 	var req repo.UserRequest
//
// 	err := c.Bind().JSON(&req)
//
// 	if err != nil {
// 		return fiber.NewError(400, "Nieprawidlowe żądanie")
// 	}
//
// 	res, err := s.DB.User().Accept(req)
//
// 	if err != nil {
// 		return err
// 	}
//
// 	return c.JSON(res)
// }
