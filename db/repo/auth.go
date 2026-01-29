package repo

import (
	"context"
	"errors"
	"log/slog"

	"github.com/MarkSmersh/nil-chat/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgconn"
)

type Auth struct {
	Repo
}

type UserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// returns a new signed access token that is generated after user's login
func (a Auth) Login(req UserLogin) (string, error) {
	hashedPassword := utils.HashPassword(req.Password)

	row := a.Pool.QueryRow(
		context.Background(),
		`select id from users where username = $1 and password = $2`,
		req.Username,
		hashedPassword,
	)

	var userId int

	if err := row.Scan(&userId); err != nil {
		return "", fiber.NewError(400, "No user's found with given credentials")
	}

	refreshToken, err := a.registerRefreshToken(userId)

	if err != nil {
		return "", err
	}

	accessToken, err := utils.NewAccessToken(refreshToken)

	if err != nil {
		return "", err
	}

	return accessToken, nil
}

type UserRegister struct {
	DisplayName string `json:"displayName"`
	Username    string `json:"username"`
	Password    string `json:"password"`
}

// returns a new signed access token that is generated after user's registation
func (a Auth) Register(req UserRegister) (string, error) {
	hashedPassword := utils.HashPassword(req.Password)

	row := a.Pool.QueryRow(
		context.Background(),
		`insert into users (username, password, display_name)
		values ($1, $2, $3) 
		returning id`,
		req.Username,
		hashedPassword,
		req.DisplayName,
	)

	var userId int

	if err := row.Scan(&userId); err != nil {
		var pgerr *pgconn.PgError

		if errors.As(err, &pgerr) {
			if pgerr.ConstraintName == "different_username" {
				return "", fiber.NewError(400, "Username's already taken")
			}

			slog.Error(err.Error())
			return "", fiber.ErrInternalServerError
		}

		return "", fiber.ErrInternalServerError
	}

	refreshToken, err := a.registerRefreshToken(userId)

	if err != nil {
		return "", err
	}

	accessToken, err := utils.NewAccessToken(refreshToken)

	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (a Auth) RevokeRefreshToken(userId int, iss string) error {
	_, err := a.Pool.Exec(
		context.Background(),
		`DELETE from users_tokens where user_id = $1 and iss = $2`,
		userId,
		iss,
	)

	if err != nil {
		slog.Error(err.Error())
		return fiber.ErrInternalServerError
	}

	return nil
}

func (a Auth) UpdateAccessToken(userId int, iss string) (string, error) {
	row := a.Pool.QueryRow(
		context.Background(),
		`SELECT refresh_token from users_tokens where iss = $1 and user_id = $2`,
		iss,
		userId,
	)

	var refreshToken string

	if err := row.Scan(&refreshToken); err != nil {
		return "", err
	}

	accessToken, err := utils.NewAccessToken(refreshToken)

	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (a Auth) registerRefreshToken(userId int) (string, error) {
	refreshToken, iss := utils.NewRefreshToken(userId)

	_, err := a.Pool.Exec(
		context.Background(),
		`insert into users_tokens (user_id, iss, refresh_token) values ($1, $2, $3)`,
		userId,
		iss,
		refreshToken,
	)

	if err != nil {
		slog.Error(err.Error())
		return "", fiber.ErrInternalServerError
	}

	return refreshToken, nil
}
