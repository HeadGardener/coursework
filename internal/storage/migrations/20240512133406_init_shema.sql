-- +goose Up
-- +goose StatementBegin
create table users (
    id uuid primary key,
    username varchar(255) not null unique,
    name varchar(255) not null,
    role integer not null,
    age integer not null,
    password_hash varchar(255) not null
);

create table drinks (
	id serial primary key,
	name varchar(255) not null unique,
	type varchar(255) not null,
    bottle integer not null default 1000,
	cost float not null,
    is_soft bool not null
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table users;
drop table drinks;
-- +goose StatementEnd
