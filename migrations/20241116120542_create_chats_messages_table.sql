-- +goose Up
-- +goose StatementBegin
create table chats
(
    id         serial primary key,
    created_at timestamp not null default now()
);

create table chat_users
(
    id      serial primary key,
    chat_id int not null,
    user_id int not null,
    foreign key (chat_id) references chats (id) ON DELETE CASCADE
);

create table messages
(
    id         serial primary key,
    chat_id    int       not null,
    user_id    int       not null,
    msg_text   text      not null,
    created_at timestamp not null default now(),
    foreign key (chat_id) references chats (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table messages;
drop table chat_users;
drop table chats;
-- +goose StatementEnd