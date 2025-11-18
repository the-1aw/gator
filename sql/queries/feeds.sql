-- name: CreateFeed :one
insert into feeds (id, created_at, updated_at, name, url, user_id)
values ($1,$2,$3,$4,$5,$6)
returning *;

-- name: GetFeedSummary :many
select feeds.name, users.name as created_by, url from feeds
left join users
on users.id = user_id;

-- name: MarkFeedFetched :exec
update feeds
set last_fetched_at = current_timestamp, updated_at = current_date
where id = $1;

-- name: GetNextFeedToFetch :one
select * from feeds
order by last_fetched_at asc nulls first
limit 1;

-- name: GetFeedByUrl :one
select * from feeds
where url = $1;

