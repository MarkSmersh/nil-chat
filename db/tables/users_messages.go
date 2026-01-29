package tables

var UsersMessagesTable = `create table if not exists users_messages (
	id bigserial primary key,
	user_id bigint references users(id) not null,
	message_id bigint references messages(id) on delete cascade not null,
	deleted_by_user boolean default false,
	unique(message_id, user_id)
)`
