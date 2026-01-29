package api

import (
	"fmt"

	"github.com/MarkSmersh/nil-chat/db"
	"github.com/gofiber/fiber/v3"
)

type Server struct {
	Api fiber.Router
	App *fiber.App
	DB  *db.DB
}

func NewServer() Server {
	app := fiber.New()

	return Server{
		App: app,
		Api: app.Group("/api"),
	}
}

func (s *Server) ConnectDB(url string) error {
	var err error
	s.DB, err = db.NewDB(url)

	if err != nil {
		return err
	}

	return err
}

func (s Server) Start(address string, port int) error {
	return s.App.Listen(fmt.Sprintf("%s:%d", address, port))
}
