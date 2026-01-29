package tables

var GlobalChatsTable = `create table if not exists global_chats (
	id bigserial primary key,
	chat_id bigint references chats(id) not null,
	title text not null
)`
