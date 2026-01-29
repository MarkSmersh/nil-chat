package tables

var MessagesTable = `create table if not exists messages (
	id bigserial primary key,
	user_id bigint references users(id) not null,

	chat_message_id bigint,
	chat_id bigint references chats(id) not null,

	text text,
	created_at timestamp default now(),
	unique(chat_message_id, chat_id)
)`
