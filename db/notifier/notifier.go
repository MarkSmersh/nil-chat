package notifier

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"slices"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Sub struct {
	UUID     string
	Callback func(Channel, Notification[json.RawMessage])
}

type Notification[T any] struct {
	UserIDs []int `json:"userIds"`
	Payload T     `json:"payload"`
}

type Notifier struct {
	conn *pgx.Conn
	pool *pgxpool.Pool
	subs []Sub

	isActive bool
}

func NewNotifier(conn *pgx.Conn, pool *pgxpool.Pool) *Notifier {
	return &Notifier{
		conn:     conn,
		pool:     pool,
		subs:     []Sub{},
		isActive: false,
	}
}

// inits the notifier and starts to listen for upcoming events and send them to the sub functions
func (n *Notifier) Listen(channels ...Channel) error {
	if !n.isActive {
		n.isActive = true
	} else {
		return errors.New("Notifier already listens")
	}

	for _, c := range channels {
		_, err := n.conn.Exec(
			context.Background(),
			fmt.Sprintf("listen %s", c),
		)

		if err != nil {
			slog.Error(err.Error())
			return err
		}
	}

	for {
		notify, err := n.conn.WaitForNotification(context.Background())

		if err != nil {
			slog.Error(err.Error())
			continue
		}

		var not Notification[json.RawMessage]

		r := bytes.NewReader([]byte(notify.Payload))

		d := json.NewDecoder(r)

		d.UseNumber()

		err = d.Decode(&not)

		if err != nil {
			slog.Error(err.Error())
			continue
		}

		for _, sub := range n.subs {
			sub.Callback(Channel(notify.Channel), not)
		}
	}
}

func (n *Notifier) Subscribe(callback func(Channel, Notification[json.RawMessage])) string {
	uuid := uuid.NewString()

	n.subs = append(n.subs, Sub{
		UUID:     uuid,
		Callback: callback,
	})

	return uuid
}

func (n *Notifier) Unsubscribe(subUUID string) {
	n.subs = slices.DeleteFunc(n.subs, func(s Sub) bool {
		return s.UUID == subUUID
	})
}

// notifies with the given channel and payload db's listeners;
// sql string must return json-like data, so all of the listeners could
// easily bind the payload data
func (n Notifier) Notify(channel Channel, data any, userIds []int) error {
	not := Notification[any]{
		UserIDs: userIds,
		Payload: data,
	}

	bytes, err := json.Marshal(not)

	if err != nil {
		slog.Error(err.Error())
		return err
	}

	_, err = n.pool.Exec(
		context.Background(),
		`select pg_notify($1, $2)`,
		channel,
		string(bytes),
	)

	if err != nil {
		slog.Error(err.Error())
		return err
	}

	return nil
}
