package tables

var UsersUsersTable = `create table if not exists users_users (
	id bigserial primary key,
	user_id bigint references users(id) not null,
	target_id bigint references users(id) not null,
	blocked boolean default false,
	constraint self_relation check (user_id != target_id),
	constraint existing_relation unique (user_id, target_id)
)`
