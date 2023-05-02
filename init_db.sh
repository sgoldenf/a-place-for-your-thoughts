#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
	CREATE USER $APP_DB_USER;
	CREATE DATABASE $APP_DB;
    ALTER USER $APP_DB_USER WITH PASSWORD '$APP_DB_PASSWORD';
    CREATE USER $TEST_DB_USER;
	CREATE DATABASE $TEST_DB;
	GRANT ALL PRIVILEGES ON DATABASE $TEST_DB TO $TEST_DB_USER;
    ALTER USER $TEST_DB_USER WITH PASSWORD '$TEST_DB_PASSWORD';

EOSQL

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$APP_DB" <<-EOSQL
	create table if not exists posts (
		id serial primary key,
		title varchar(255) not null,
		text text not null
	);

	create table if not exists sessions (
		token text primary key,
		data bytea not null,
		expiry timestamptz not null
	);
	create index sessions_expiry_idx on sessions (expiry);

	create table if not exists users (
		id serial primary key,
		name varchar(255) not null,
		email varchar(255) not null,
		hashed_password char(60) not null,
		created timestamp without time zone default (now() at time zone 'utc')
	);
	alter table users add constraint users_uc_email UNIQUE(email);
	
	GRANT CONNECT ON DATABASE $APP_DB TO $APP_DB_USER;
	GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO $APP_DB_USER;
EOSQL
