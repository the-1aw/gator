-- +goose Up
-- +goose StatementBegin
create table posts(
	id uuid default gen_random_uuid() primary key,
	created_at timestamp default current_timestamp not null,
	updated_at timestamp default current_timestamp not null,
	title text not null,
	url text unique not null,
	description text,
	published_at timestamp,
	feed_id uuid not null references feeds(id) on delete cascade
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table posts;
-- +goose StatementEnd
