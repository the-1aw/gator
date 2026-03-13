-- name: CreateFeedFollow :many
with inserted_feed_follow as (
insert into feed_follows (user_id, feed_id)
values ($1, $2)
on conflict(user_id, feed_id) do update set updated_at = current_timestamp
returning *
)
select inserted_feed_follow.*, users.name as username, feeds.name as feed_name from inserted_feed_follow
inner join users on users.id = inserted_feed_follow.user_id
inner join feeds on feeds.id = inserted_feed_follow.feed_id;

-- name: DeleteFeedFollow :exec
delete from feed_follows
where user_id = $1 and feed_id = $2;
