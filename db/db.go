package db

import (
	"context"

	"github.com/MarkSmersh/nil-chat/db/notifier"
	"github.com/MarkSmersh/nil-chat/db/repo"
	"github.com/MarkSmersh/nil-chat/db/tables"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	repo   repo.Repo
	notify *notifier.Notifier
	pool   *pgxpool.Pool
}

func NewDB(url string) (*DB, error) {
	var db DB

	pool, err := pgxpool.New(context.Background(), url)

	if err != nil {
		return nil, err
	}

	tables.Init(pool)

	db.pool = pool

	conn, err := pool.Acquire(context.Background())

	if err != nil {
		return nil, err
	}

	db.notify = notifier.NewNotifier(conn.Conn(), pool)

	go db.notify.Listen(notifier.DeleteMessage, notifier.NewMessage)

	db.repo = repo.NewRepo(pool, db.notify)

	return &db, nil
}

func (db DB) WithAuth(userId int) *DB {
	db.repo.SetUserID(userId)
	return &db
}

func (db DB) User() repo.User {
	return repo.User{Repo: db.repo}
}

func (db DB) Message() repo.Message {
	return repo.Message{Repo: db.repo}
}

func (db DB) Chat() repo.Chat {
	return repo.Chat{Repo: db.repo}
}

func (db DB) Auth() repo.Auth {
	return repo.Auth{Repo: db.repo}
}

func (db *DB) Notifier() *notifier.Notifier {
	return db.notify
}
