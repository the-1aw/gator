-- +goose Up
-- +goose StatementBegin
create table users (
	id uuid primary key,
	created_at timestamp,
	updated_at timestamp,
	name text not null,
	unique(name)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table users;
-- +goose StatementEnd
