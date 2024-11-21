-- +goose Up
-- +goose StatementBegin
create table chats (
    id serial primary key,
    created_at timestamp not null default now(),
    updated_at timestamp
);
create table chat_users (
   chat_id int not null,
   user_id int not null,
   primary key (chat_id, user_id),
   foreign key (chat_id) references chats(id)
);
create table messages (
    id serial primary key,
    chat_id int not null,
    from_user_id int not null,
    msg_txt text not null,
    created_at timestamp not null default now(),
    updated_at timestamp,
    foreign key (chat_id) references chats(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table messages;
drop table chat_users;
drop table chats;
-- +goose StatementEnd
