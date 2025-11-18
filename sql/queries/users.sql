-- name: CreateUser :one
insert into users (id, created_at, updated_at, name)
values ($1, $2, $3, $4)
returning *;

-- name: GetUser :one
select * from users
where name = $1;

-- name: GetUsers :many
select * from users;

-- name: GetUserFeedFollow :many
select users.name as username, feeds.name as feed_name from users
join feed_follows on feed_follows.user_id = users.id
join feeds on feed_follows.feed_id = feeds.id
where users.name = $1;

-- name: DeleteAllUsers :exec
delete from users;
