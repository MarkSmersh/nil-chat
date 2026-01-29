package tables

// put password and etc to the account table?

var UsersTable = `create table if not exists users (
	id bigserial primary key,
	username text constraint different_username unique not null,
	password text not null,
	display_name text not null,
	last_seen_at timestamp default now(),
	created_at timestamp default now()
)`
