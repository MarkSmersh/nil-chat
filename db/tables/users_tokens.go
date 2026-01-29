package tables

var UsersTokensTable = `create table if not exists users_tokens (
	id bigserial primary key,
	iss uuid not null,
	user_id bigint references users(id) not null,
	refresh_token text not null,
	unique(user_id, iss)
)`
