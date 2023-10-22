-- +goose Up
-- +goose StatementBegin
create table chat (
    id serial primary key,
    name varchar(64)
);

create table chats_to_users (
    chatId int references chat(id) on delete cascade,
    username varchar(64),

    unique(chatId, username)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table chats_to_users;
drop table chat;
-- +goose StatementEnd
