package tables

var ChatsTable = `create table if not exists chats (
	id bigserial primary key,
	type chat_type not null
)`
