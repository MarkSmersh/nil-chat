package tables

var PrivateChatsTable = `create table if not exists private_chats (
	id bigserial primary key,
	chat_id bigint references chats(id),
	a_user_id bigint references users(id),
	b_user_id bigint references users(id),
	constraint invalid_insert check (a_user_id <= b_user_id),
	unique (chat_id, a_user_id, b_user_id)
)`
