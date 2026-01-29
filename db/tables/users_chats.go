package tables

var UsersChatsTable = `create table if not exists users_chats (
	id bigserial primary key,
	user_id bigint references users(id) not null,
	chat_id bigint references chats(id) not null,
	unique (user_id, chat_id)
)`
