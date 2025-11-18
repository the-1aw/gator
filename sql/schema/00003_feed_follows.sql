-- +goose Up
-- +goose StatementBegin
create table feed_follows(
id uuid primary key,
	created_at date not null,
	updated_at date,
	user_id uuid not null references users(id) on delete cascade,
	feed_id uuid not null references feeds(id) on delete cascade,
	unique(user_id, feed_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table feed_follows
-- +goose StatementEnd
