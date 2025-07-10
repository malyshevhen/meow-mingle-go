-- Create keyspace with proper replication;

CREATE KEYSPACE IF NOT EXISTS mingle WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };

-- User profiles table;

CREATE TABLE IF NOT EXISTS mingle.profiles (
    user_id text PRIMARY KEY,
    email text,
    first_name text,
    last_name text,
    bio text,
    avatar_url text,
    created_at timestamp,
    updated_at timestamp
);

-- Create index for email lookups;

CREATE INDEX IF NOT EXISTS profiles_email_idx ON mingle.profiles (email);

-- Posts table with time-series optimization;

CREATE TABLE IF NOT EXISTS mingle.posts (
    id uuid PRIMARY KEY,
    author_id text,
    content text,
    image_urls list<text>,
    created_at timestamp,
    updated_at timestamp
);

-- Posts by author table for efficient user timeline queries;

CREATE TABLE IF NOT EXISTS mingle.posts_by_author (
    author_id text,
    created_at timestamp,
    post_id uuid,
    content text,
    image_urls list<text>,
    updated_at timestamp,
PRIMARY KEY (author_id, created_at, post_id)
) WITH CLUSTERING ORDER BY (created_at DESC, post_id ASC);

-- Comments table;

CREATE TABLE IF NOT EXISTS mingle.comments (
    id uuid PRIMARY KEY,
    post_id uuid,
    author_id text,
    content text,
    created_at timestamp,
    updated_at timestamp
);

-- Comments by post table for efficient post comment queries;

CREATE TABLE IF NOT EXISTS mingle.comments_by_post (
    post_id uuid,
    created_at timestamp,
    comment_id uuid,
    author_id text,
    content text,
    updated_at timestamp,
PRIMARY KEY (post_id, created_at, comment_id)
) WITH CLUSTERING ORDER BY (created_at DESC, comment_id ASC);

-- Reactions table with proper composite key;

CREATE TABLE IF NOT EXISTS mingle.reactions (
    target_id uuid,
    target_type text,
    author_id text,
    reaction_type text,
    created_at timestamp,
PRIMARY KEY ((target_id, target_type), author_id)
);

-- Reactions by target table for efficient reaction counting;

CREATE TABLE IF NOT EXISTS mingle.reactions_by_target (
    target_id uuid,
    target_type text,
    reaction_type text,
    author_id text,
    created_at timestamp,
PRIMARY KEY ((target_id, target_type), reaction_type, author_id)
);

-- Subscriptions/Following table;

CREATE TABLE IF NOT EXISTS mingle.subscriptions (
    follower_id text,
    following_id text,
    created_at timestamp,
PRIMARY KEY (follower_id, following_id)
);

-- Followers table for reverse lookup;

CREATE TABLE IF NOT EXISTS mingle.followers (
    following_id text,
    follower_id text,
    created_at timestamp,
PRIMARY KEY (following_id, follower_id)
);

-- User feed table for timeline generation;

CREATE TABLE IF NOT EXISTS mingle.user_feed (
    user_id text,
    created_at timestamp,
    post_id uuid,
    author_id text,
    content text,
    image_urls list<text>,
PRIMARY KEY (user_id, created_at, post_id)
) WITH CLUSTERING ORDER BY (created_at DESC, post_id ASC);

-- User activity log for audit trail;

CREATE TABLE IF NOT EXISTS mingle.user_activity (
    user_id text,
    activity_id uuid,
    activity_type text,
    target_id uuid,
    metadata map<text, text>,
    created_at timestamp,
PRIMARY KEY (user_id, created_at, activity_id)
) WITH CLUSTERING ORDER BY (created_at DESC, activity_id ASC);
