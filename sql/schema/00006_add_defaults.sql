-- +goose Up
-- +goose StatementBegin
alter table users
    alter column id set default gen_random_uuid(),
    alter column created_at type timestamp,
    alter column created_at set default current_timestamp,
    alter column updated_at type timestamp,
    alter column updated_at set default current_timestamp;

alter table feeds
    alter column id set default gen_random_uuid(),
    alter column created_at type timestamp using created_at::timestamp,
    alter column created_at set default current_timestamp,
    alter column updated_at type timestamp using updated_at::timestamp,
    alter column updated_at set default current_timestamp;

alter table feed_follows
    alter column id set default gen_random_uuid(),
    alter column created_at type timestamp using created_at::timestamp,
    alter column created_at set default current_timestamp,
    alter column updated_at type timestamp using updated_at::timestamp,
    alter column updated_at set default current_timestamp;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table users
    alter column id drop default,
    alter column created_at drop default,
    alter column updated_at drop default;

alter table feeds
    alter column id drop default,
    alter column created_at type date using created_at::date,
    alter column created_at drop default,
    alter column updated_at type date using updated_at::date,
    alter column updated_at drop default;

alter table feed_follows
    alter column id drop default,
    alter column created_at type date using created_at::date,
    alter column created_at drop default,
    alter column updated_at type date using updated_at::date,
    alter column updated_at drop default;
-- +goose StatementEnd
