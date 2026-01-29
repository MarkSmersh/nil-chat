package repo

import (
	"context"
	"log/slog"

	"github.com/MarkSmersh/nil-chat/models"
	"github.com/gofiber/fiber/v3"
)

type User struct {
	Repo
}

func (u User) Search(prompt string) ([]models.User, error) {
	if prompt == "" {
		return []models.User{}, nil
	}

	row := u.Pool.QueryRow(
		context.Background(),
		`
SELECT
    json_agg(
        json_build_object(
            'id',
            id,
            'username',
            username,
            'createdAt',
            floor(
                extract(
                    epoch
                    FROM
                        created_at
                )
            ),
            'displayName',
            display_name,
			'lastSeenAt',
            floor(
                extract(
                    epoch
                    FROM
                        last_seen_at
                )
            )
        )
    )
FROM
    users
WHERE
    username || ' ' || display_name LIKE $1;


		`,
		"%"+prompt+"%",
	)

	users := []models.User{}

	if err := row.Scan(&users); err != nil {
		slog.Error(err.Error())
		return []models.User{}, fiber.ErrInternalServerError
	}

	if users == nil {
		return []models.User{}, nil
	}

	return users, nil
}

type UserBlock struct {
	Target int `json:"target"`
}

func (u User) Block(req UserBlock) (bool, error) {
	tag, err := u.Pool.Exec(
		context.Background(),
		`
		insert into users_users (user_id, target_id, blocked)
		values ($1, $2, true), ($2, $1, false)
		on conflict on constraint existing_relation do update user_users
		set blocked = true
		where user_id = $1
			and target_id = $2
			and blocked = false
		`,
		u.UserID,
		req.Target,
	)

	if tag.RowsAffected() <= 0 {
		return false, fiber.NewError(400, "Użytkownik już jest zablokowany")
	}

	if err != nil {
		slog.Error(err.Error())
		return false, fiber.ErrInternalServerError
	}

	return true, nil
}

type UserUnblock struct {
	TargetID int `json:"targetID"`
}

func (u User) Unblock(req UserUnblock) (bool, error) {
	tag, err := u.Pool.Exec(
		context.Background(),
		`
		update users_users set blocked = false
		where user_id = $1 
			and target_id = $2
			and blocked = true
		`,
		u.UserID,
		req.TargetID,
	)

	if tag.RowsAffected() <= 0 {
		return false, fiber.NewError(400, "Użytkownik nie jest zablokowany")
	}

	if err != nil {
		slog.Error(err.Error())
		return false, fiber.ErrInternalServerError
	}

	return true, nil
}

// type UserRequest struct {
// 	TargetID int `json:"targetID"`
// }
//
// func (u User) Request(req UserRequest) (bool, error) {
// 	tag, err := u.Pool.Exec(
// 		context.Background(),
// 		`
// 		insert into users_users (user_id, target_id, accepted, requested)
// 		values ($1, $2, true, true)
// 		on conflict
// 			on constraint existing_relation
// 			do
// 		update user_users
// 		set accepted = true,
// 			requested = true
// 		where user_id = $1
// 			and
// 			target_id = $2
// 			and requested = false
// 		`,
// 		u.UserID,
// 		req.TargetID,
// 	)
//
// 	if err != nil {
// 		slog.Error(err.Error())
// 		return false, fiber.ErrInternalServerError
// 	}
//
// 	if tag.RowsAffected() == 0 {
// 		return false, fiber.NewError(400, "You've already sent a request to this user")
// 	}
//
// 	return true, nil
// }
//
// func (u User) Accept(req UserRequest) (bool, error) {
// 	tag, err := u.Pool.Exec(
// 		context.Background(),
// 		`
// 		update user_users
// 		set accepted = true
// 		where user_id = $1
// 			and target_id = $2
// 			and requested = false
// 		`,
// 		u.UserID,
// 		req.TargetID,
// 	)
//
// 	if err != nil {
// 		slog.Error(err.Error())
// 		return false, fiber.ErrInternalServerError
// 	}
//
// 	if tag.RowsAffected() == 0 {
// 		return false, fiber.NewError(400, "You've already sent a request to this user")
// 	}
//
// 	return true, nil
// }
