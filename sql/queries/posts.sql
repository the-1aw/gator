-- name: CreatePost :one
insert into posts (title, url, description, published_at, feed_id)
values ($1, $2, $3, $4, $5)
on conflict (url) do update set title = $1, description = $3, published_at = $4, updated_at = current_timestamp
returning *;

-- name: GetPostForUser :many
select posts.* from posts
inner join feeds on posts.feed_id = feeds.id
inner join feed_follows on feed_follows.feed_id = feeds.id
where feed_follows.user_id = $1
order by posts.published_at asc nulls last
limit $2;
