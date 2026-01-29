package repo

import (
	"github.com/MarkSmersh/nil-chat/db/notifier"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	Pool     *pgxpool.Pool
	UserID   int
	Notifier *notifier.Notifier
}

func NewRepo(pool *pgxpool.Pool, notifier *notifier.Notifier) Repo {
	return Repo{
		Pool:     pool,
		Notifier: notifier,
	}
}

func (r *Repo) SetUserID(userId int) {
	r.UserID = userId
}
