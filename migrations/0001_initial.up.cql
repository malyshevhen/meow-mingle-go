
CREATE KEYSPACE IF NOT EXISTS mingle WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };

CREATE TABLE IF NOT EXISTS mingle.profiles (
    user_id text PRIMARY KEY,
    email text,
    first_name text,
    last_name text,
    created_at timestamp,
    updated_at timestamp
);

CREATE TABLE IF NOT EXISTS mingle.posts (
    id uuid PRIMARY KEY,
    author_id text,
    content text,
    created_at timestamp,
    updated_at timestamp
);

CREATE TABLE IF NOT EXISTS mingle.comments (
    id uuid PRIMARY KEY,
    post_id uuid,
    author_id text,
    content text,
    created_at timestamp,
    updated_at timestamp
);

CREATE TABLE IF NOT EXISTS mingle.reactions (
    target_id uuid,
    author_id text,
    content text,
    created_at timestamp,
    PRIMARY KEY (target_id, author_id)
);

CREATE TABLE IF NOT EXISTS mingle.subscriptions (
    follower_id text,
    following_id text,
    created_at timestamp,
    PRIMARY KEY (follower_id, following_id)
);
