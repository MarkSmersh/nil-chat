package repo

import (
	"context"
	"errors"
	"log/slog"

	"github.com/MarkSmersh/nil-chat/models"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Chat struct {
	Repo
}

func (c Chat) JoinGlobal() (models.Chat, error) {
	// FIXME: The next statement declarates, that
	// the global chat could be the only and the only one

	row := c.Pool.QueryRow(
		context.Background(),
		`
with chat as (
	insert into chats (type) select 'global'
	where not exists (
		select 1 from global_chats
	)
	returning id
),
global_chat as (
	insert into global_chats (chat_id, title)
	select id, $1 from chat
	returning *
),
uc as (
	insert into users_chats (user_id, chat_id)
	select $2, c.chat_id from (
		select chat_id
		from global_chat

		union

		select chat_id
		from global_chats
		limit 1
	) c
)
select json_build_object(
	'id', c.chat_id,
	'type', 'global',
	'title', c.title
)
from (
	select chat_id, title
	from global_chat

	union

	select chat_id, title
	from global_chats
	limit 1
) c
		`,
		uuid.NewString(),
		c.UserID,
	)

	var chat models.Chat

	if err := row.Scan(&chat); err != nil {
		slog.Error(err.Error())
		return chat, fiber.ErrInternalServerError
	}

	return chat, nil
}

type ChatGetHistory struct {
	ChatID int `json:"chatId"`
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

func (c Chat) GetHistory(req ChatGetHistory) ([]models.Message, error) {
	if req.Limit <= 0 || req.Limit > 500 {
		return []models.Message{}, fiber.NewError(400, "parametr limit nie może być mniej 0 i więcej 500")
	}

	if req.ChatID <= 0 {
		return []models.Message{}, fiber.NewError(400, "brak parametru chatId")
	}

	row := c.Pool.QueryRow(
		context.Background(),
		`
		select json_agg(
    		json_build_object(
    		    'id',
    		    m.id,
    		    'text',
    		    m.text,
    		    'fromId',
    		    m.user_id,
    		    'chatId',
    		    m.chat_id,
    		    'timestamp',
    		    floor(
    		        extract(
    		            epoch
    		            FROM
    		                m.created_at
    		        )
    		    )
    		)
		)
		from (
			select *
			from messages
			where chat_id = (
				select chat_id
				from users_chats
				where chat_id = $1
					and user_id = $2
			)
			order by id desc
			offset $3 
			limit $4
		) m
		`,
		req.ChatID,
		c.UserID,
		req.Offset,
		req.Limit,
	)

	var msgs []models.Message

	if err := row.Scan(&msgs); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []models.Message{}, nil
		}

		slog.Error(err.Error())
		return []models.Message{}, fiber.ErrInternalServerError
	}

	return msgs, nil
}

func (c Chat) GetAll() ([]models.Chat, error) {
	row := c.Pool.QueryRow(
		context.Background(),
		`
SELECT
	json_agg(
		json_build_object(
		    'id',
		    c.id,
		    'type',
		    c.type,
		    'username',
		    u.username,
		    'displayName',
		    u.display_name,
		    'title',
		    gc.title
		)
	)
FROM
    users_chats uc
    left JOIN chats c ON c.id = uc.chat_id
    left JOIN (
        SELECT
            *
        FROM
            private_chats
        WHERE
            a_user_id = $1
            OR b_user_id = $1
    ) pc ON pc.chat_id = uc.chat_id
    JOIN lateral (
        SELECT
            *
        FROM
            users u
        WHERE
            u.id IN (pc.a_user_id, pc.b_user_id)
    ) u ON TRUE
    LEFT JOIN global_chats gc ON gc.chat_id = uc.chat_id
WHERE
    user_id = $1;

		`,
		c.UserID,
	)

	var chats []models.Chat

	if err := row.Scan(&chats); err != nil {
		slog.Error(err.Error())
		return []models.Chat{}, nil
	}

	return chats, nil
}

type ChatFromUser struct {
	UserID int `json:"userId"`
}

func (c Chat) FromUser(req ChatFromUser) (models.Chat, error) {
	if req.UserID == 0 {
		return models.Chat{}, fiber.NewError(400, "Brak parametry userId")
	}

	row := c.Pool.QueryRow(
		context.Background(),
		`
		WITH chat AS (
    INSERT INTO
        chats (TYPE)
    SELECT
        'private'
    WHERE
        NOT EXISTS (
            SELECT
                chat_id
            FROM
                private_chats
            WHERE
                a_user_id = least($1::bigint, $2::bigint)::bigint
                AND b_user_id = greatest($1::bigint, $2::bigint)::bigint
        )
    RETURNING
		*
),
private_chat AS (
    INSERT INTO
        private_chats (chat_id, a_user_id, b_user_id)
    SELECT
        id,
        least($1, $2)::bigint,
        greatest($1, $2)::bigint
    FROM
        chat ON conflict DO nothing
    RETURNING
        *
),
user_user AS (
    INSERT INTO
        users_users (user_id, target_id)
    SELECT
        u.a,
        u.b
    FROM
        (
            VALUES
                ($1::bigint, $2::bigint),
                ($2::bigint, $1::bigint)
        ) AS u (a, b)
    WHERE
        NOT EXISTS (
            SELECT
                1
            FROM
                users_users
            WHERE
                (
                    user_id = $1::bigint
                    AND target_id = $2::bigint
                )
                OR (
                    user_id = $2::bigint
                    AND target_id = $1::bigint
                )
        )
        AND $1::bigint != $2::bigint
)
SELECT
    json_build_object(
        'id', c.id,
        'type', 'private',
        'displayName', c.display_name,
        'username', c.username
    )
FROM (
    SELECT
        c.id,
        u.display_name,
        u.username
    FROM chat c
    JOIN users u ON u.id = $1

    UNION

    SELECT
        pc.id,
        u.display_name,
        u.username
    FROM private_chats pc
    JOIN users u ON u.id = $2
    WHERE
        pc.a_user_id = least($1::bigint, $2::bigint)
        AND pc.b_user_id = greatest($1::bigint, $2::bigint)
) c;
		`,
		req.UserID,
		c.UserID,
	)

	var chat models.Chat

	if err := row.Scan(&chat); err != nil {
		var pgerr *pgconn.PgError

		if errors.As(err, &pgerr) {
			if pgerr.ConstraintName == "users_users_user_id_fkey" ||
				pgerr.ConstraintName == "private_chats_b_user_id_fkey" {
				return models.Chat{}, fiber.NewError(404, "Nie ma użytkownika z podanym userId")
			}
		}

		slog.Error(err.Error())
		return models.Chat{}, fiber.ErrInternalServerError
	}

	return chat, nil
}

func (c Chat) GetPreviews() ([]models.Message, error) {
	row := c.Pool.QueryRow(
		context.Background(),
		`
SELECT
    json_agg(
        json_build_object(
            'id',
            m.chat_message_id,
            'chatId',
            m.chat_id,
            'fromId',
            m.user_id,
            'text',
            m.text,
            'timestamp',
            floor(
                extract (
                    epoch
                    FROM
                        m.created_at
                )
            )
        )
    )
FROM
    users_chats uc
    LEFT JOIN LATERAL (
        SELECT
            *
        FROM
            messages m
        WHERE
            m.id not in (
                SELECT
                    message_id
                FROM
                    users_messages
                WHERE
                    deleted_by_user = TRUE
					and user_id = $1
            )
            AND m.chat_id = uc.chat_id
        ORDER BY
            ID DESC
        LIMIT
            1
    ) m ON TRUE
WHERE
    uc.user_id = $1;
		`,
		c.UserID,
	)

	var msgs []models.Message

	if err := row.Scan(&msgs); err != nil {
		slog.Error(err.Error())
		return []models.Message{}, fiber.ErrInternalServerError
	}

	return msgs, nil
}
