create table if not exists posts (
    id serial primary key,
    title varchar(255) not null,
    text text not null
);
create table if not exists users (
    id serial primary key,
    name varchar(255) not null,
    email varchar(255) not null,
    hashed_password char(60) not null,
    created timestamp without time zone default (now() at time zone 'utc')
);
alter table users
add constraint users_uc_email UNIQUE(email);
insert into users (name, email, hashed_password, created)
values (
        'sgoldenf',
        'sgoldenf@example.com',
        '$2a$12$NuTjWXm3KKntReFwyBVHyuf/to.HEwTy.eS206TNfkGfr6HzGJSWG',
        '2022-01-01 10:00:00'
    );
