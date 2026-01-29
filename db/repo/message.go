package repo

import (
	"context"
	"errors"
	"log/slog"

	"github.com/MarkSmersh/nil-chat/db/notifier"
	"github.com/MarkSmersh/nil-chat/models"
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Message struct {
	Repo
}

type MessageSend struct {
	ChatID int    `json:"chatId"`
	Text   string `json:"text"`
}

func (m Message) Send(req MessageSend) (models.Message, error) {
	if req.ChatID == 0 {
		return models.Message{}, fiber.NewError(400, "Brak parametru chatId")
	}

	if req.Text == "" {
		return models.Message{}, fiber.NewError(400, "Brak parametru text")
	}

	row := m.Pool.QueryRow(
		context.Background(),
		`
WITH pc AS (
    SELECT
        a_user_id,
        b_user_id
    FROM
        private_chats
    WHERE
        chat_id = $3
        AND (
            a_user_id = $1
            OR b_user_id = $1
        )
	limit 1
),
chat_a AS (
    INSERT INTO
        users_chats (user_id, chat_id)
	select unnest(array[a_user_id, b_user_id]), $3 from pc
    ON conflict DO nothing
),
last_message_in_chat AS (
    SELECT
        chat_message_id
    FROM
        messages
    WHERE
        chat_id = $3
    ORDER BY
        id DESC
    LIMIT
        1
), blocked_by_user AS (
    SELECT
        $3
    FROM
        users_users
    WHERE
        blocked = TRUE
        AND target_id = $1
        AND (
            user_id = (
                SELECT
                    pc.a_user_id
                FROM
                    pc
            )
            OR user_id = (
                SELECT
                    pc.b_user_id
                FROM
                    pc
            )
        )
)
INSERT INTO
    messages (user_id, text, chat_id, chat_message_id)
SELECT
    $1,
    $2,
    $3,
    coalesce(
        (
            SELECT
                chat_message_id
            FROM
                last_message_in_chat
        ),
        0
    ) + $3
WHERE
    EXISTS (
        SELECT
            $3
        FROM
            pc
    )
    AND NOT EXISTS (
        SELECT
            $3
        FROM
            blocked_by_user
    )
RETURNING
    json_build_object(
        'id',
        id,
        'text',
        text,
        'fromId',
        user_id,
        'chatId',
        chat_id,
        'timestamp',
        floor(
            extract(
                epoch
                FROM
                    created_at
            )
        )
    );


		`,
		m.UserID,
		req.Text,
		req.ChatID,
	)

	var msg models.Message

	if err := row.Scan(&msg); err != nil {
		var pgerr *pgconn.PgError

		if errors.As(err, &pgerr) {
			if pgerr.ConstraintName == "chats_messages_chat_id_fkey" ||
				pgerr.ConstraintName == "messages_chat_id_fkey" {
				return msg, fiber.NewError(400, "Such chat id doesn't exist")
			}

			slog.Error(err.Error())
			return msg, fiber.ErrInternalServerError
		}

		if errors.Is(err, pgx.ErrNoRows) {
			return msg, fiber.NewError(403, "Nie masz dostępu do tego czatu")
		}

		slog.Error(err.Error())
		return msg, fiber.ErrInternalServerError
	}

	go func() {
		var userIds []int

		row := m.Pool.QueryRow(
			context.Background(),
			`
			select array_agg(user_id)
			from users_chats
				where chat_id = $1
					and user_id != $2
			`,
			req.ChatID,
			m.UserID,
		)

		if err := row.Scan(&userIds); err != nil {
			slog.Error(err.Error())
			return
		}

		m.Notifier.Notify(
			notifier.NewMessage,
			msg,
			userIds,
		)
	}()

	return msg, nil
}

type MessageDelete struct {
	ChatID    int  `json:"chatId"`
	MessageID int  `json:"messageId"`
	ForAll    bool `json:"forAll"`
}

func (m Message) Delete(req MessageDelete) (bool, error) {
	if req.ChatID == 0 {
		return false, fiber.NewError(400, "Brak parametru chatId")
	}

	if req.MessageID == 0 {
		return false, fiber.NewError(400, "Brak parametru messageId")
	}

	var (
		row pgx.Row
	)

	if req.ForAll {
		row = m.Pool.QueryRow(
			context.Background(),
			`
			delete from messages
			where chat_message_id = $1 
				and chat_id = $2
				and id not in (
					select message_id
					from users_messages
					where deleted_by_user = true
						and user_id = $3
				)
				and exists (
					select 1
					from users_chats
					where user_id = $3
						and chat_id = $2
				)
			returning json_build_object(
				'id', $1,
				'chatId', $2
			)
			`,
			req.MessageID,
			req.ChatID,
			m.UserID,
		)
	} else {
		row = m.Pool.QueryRow(
			context.Background(),
			`
			insert into users_messages um
				(message_id, deleted_by_user, user_id)
			select id, true, $1
			from messages
			where chat_message_id = $3
				and chat_id = $2
				and exists (
					select 1
					from users_chats
					where user_id = $1
						and chat_id = $2
				)
				limit 1
			returning json_build_object(
				'id', $3,
				'chatId', $2
			)
			`,
			m.UserID,
			req.ChatID,
			req.MessageID,
		)
	}

	var msg models.DeleteMessage

	if err := row.Scan(&msg); err != nil {
		slog.Error(err.Error())

		var pgerr *pgconn.PgError
		if errors.As(err, &pgerr) {
			if pgerr.ConstraintName == "users_messages_message_id_user_id_key" {
				return false, fiber.NewError(400, "Wiadomość już usunięto")
			}
		}

		if errors.Is(err, pgx.ErrNoRows) {
			return false, fiber.NewError(400, "Wiadomości z podanym id nie istnieje")
		}

		slog.Error(err.Error())
		return false, fiber.ErrInternalServerError
	}

	if req.ForAll {
		go func() {
			var userIds []int

			row := m.Pool.QueryRow(
				context.Background(),
				`
			select array_agg(user_id)
			from users_chats
				where chat_id = $1
					and user_id != $2
			`,
				req.ChatID,
				m.UserID,
			)

			if err := row.Scan(&userIds); err != nil {
				slog.Error(err.Error())
				return
			}

			m.Notifier.Notify(
				notifier.DeleteMessage,
				msg,
				userIds,
			)
		}()
	}

	return true, nil
}
