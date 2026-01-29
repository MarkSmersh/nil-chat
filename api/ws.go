package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"slices"

	"github.com/MarkSmersh/nil-chat/api/service"
	"github.com/MarkSmersh/nil-chat/db/notifier"
	"github.com/gofiber/contrib/v3/websocket"
)

func (s *Server) WSConnection(c *websocket.Conn) {
	userId := c.Locals("userId").(int)

	var (
		data []byte
		err  error
	)

	chats, err := s.DB.WithAuth(userId).Chat().GetAll()

	if err != nil {
		if err := c.WriteJSON(service.NewError("", err.Error())); err != nil {
			slog.Error(err.Error())
			return
		}
	}

	uuid := s.DB.Notifier().Subscribe(func(ch notifier.Channel, n notifier.Notification[json.RawMessage]) {
		slog.Info(
			fmt.Sprintf("%v", chats),
		)

		_, res := service.NewUpdate(
			ch,
			n.Payload,
		)

		if slices.Contains(n.UserIDs, userId) {
			if err := c.WriteJSON(res); err != nil {
				slog.Error(err.Error())
				return
			}
		}
	})

	defer s.DB.Notifier().Unsubscribe(uuid)

	for {
		if _, data, err = c.ReadMessage(); err != nil {
			slog.Error(err.Error())
			break
		}

		res := service.NewRequest(s.DB).
			WithAuth(userId).
			Send(data)

		if err := c.WriteJSON(res); err != nil {
			slog.Error(err.Error())
			break
		}
	}
}
