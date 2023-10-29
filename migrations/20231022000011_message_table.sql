-- +goose Up
-- +goose StatementBegin
create table message (
    id serial primary key,
    "from" varchar(64),
    "text" text,
    created_at timestamp default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table message;
-- +goose StatementEnd
