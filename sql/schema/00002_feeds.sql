-- +goose Up
-- +goose StatementBegin
create table feeds (
	id uuid primary key,
	created_at date,
	updated_at date,
	name text not null,
	url text not null,
	user_id uuid not null references users(id) on delete cascade,
	unique(url)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table feeds;
-- +goose StatementEnd
