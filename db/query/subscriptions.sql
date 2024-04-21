-- name: CreateSubscription :exec
INSERT INTO users_subscriptions (
    user_id, subscription_id
) VALUES ($1, $2);

-- name: GetSubscription :one
SELECT * FROM users_subscriptions
WHERE user_id = $1 AND subscription_id = $2
LIMIT 1;

-- name: ListSubscriptions :many
SELECT subscription_id
FROM users_subscriptions
WHERE user_id = $1;

-- name: DeleteSubscription :exec
DELETE FROM users_subscriptions
WHERE user_id = $1 AND subscription_id = $2;
