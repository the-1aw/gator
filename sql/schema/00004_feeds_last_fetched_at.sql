-- +goose Up
-- +goose StatementBegin
alter table feeds add last_fetched_at timestamp;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table feeds drop column last_fetched_at;
-- +goose StatementEnd
