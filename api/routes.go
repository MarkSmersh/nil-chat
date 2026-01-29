package api

import (
	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
)

func (s *Server) InitRoutes() {
	s.HealthyRoute()
	s.WSRoute()
	s.AuthRoute()
	s.ChatRoute()
	s.UserRoute()
}

func (s *Server) HealthyRoute() {
	s.Api.Get("/healthy", s.AuthMiddleware, s.HealthyGet)
}

func (s *Server) WSRoute() {
	ws := s.App.Group("/ws")

	ws.Use(func(c fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}

		return fiber.ErrUpgradeRequired
	})

	ws.Get("/:id", s.AuthMiddleware, websocket.New(s.WSConnection))
}

func (s *Server) AuthRoute() {
	g := s.Api.Group("/auth")

	g.Get("/logout", s.AuthLogoutGet)

	g.Post("/login", s.AuthLoginPost)
	g.Post("/register", s.AuthRegisterPost)
}

func (s *Server) ChatRoute() {
	g := s.Api.Group("/chat", s.AuthMiddleware)

	g.Get("/getPreviews", s.ChatGetPreviews)
	g.Get("/getAll", s.ChatGetAll)
	g.Get("/joinGlobal", s.ChatJoinGlobal)

	g.Post("/fromUser", s.ChatFromUser)
	g.Post("/getHistory", s.ChatGetHistory)
}

func (s *Server) UserRoute() {
	g := s.Api.Group("/user", s.AuthMiddleware)

	g.Get("/search/:q", s.UserSearch)

	g.Post("/block", s.UserBlock)
	g.Post("/unblock", s.UserUnblock)
}
