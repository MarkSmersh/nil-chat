package tables

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var tables = []string{
	Enums,

	ChatsTable,
	UsersTable,
	MessagesTable,

	UsersTokensTable,
	UsersUsersTable,
	UsersMessagesTable,
	UsersChatsTable,

	GlobalChatsTable,
	PrivateChatsTable,
}

func Init(pool *pgxpool.Pool) {
	for _, t := range tables {
		createTable(pool, t)
	}
}
func createTable(pool *pgxpool.Pool, query string) {
	_, err := pool.Exec(context.Background(), query)

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		var pgerr *pgconn.PgError
		if errors.As(err, &pgerr) {
			if pgerr.Code != "42710" {
				slog.Error(
					fmt.Sprintf(
						"While trying to create a scheme error occured. Query: %s",
						query,
					),
				)
				slog.Error(err.Error())
				os.Exit(1)
			}

		} else {
			slog.Error(err.Error())
			os.Exit(1)
		}
	}
}
