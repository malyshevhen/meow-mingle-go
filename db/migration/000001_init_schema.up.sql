CREATE TABLE IF NOT EXISTS users (
  id         BIGSERIAL    NOT NULL,
  email      VARCHAR(255) NOT NULL CHECK (length(email) > 0) UNIQUE,
  first_name VARCHAR(255) NOT NULL CHECK (length(first_name) > 0),
  last_name  VARCHAR(255) NOT NULL CHECK (length(last_name) > 0),
  password   VARCHAR(255) NOT NULL CHECK (length(password) > 0),
  created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS posts (
  id         BIGSERIAL    NOT NULL,
  content    VARCHAR(255) NOT NULL CHECK (length(content) > 0),
  author_id  BIGINT       NOT NULL,
  created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

  PRIMARY KEY (id),
  FOREIGN KEY (author_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS comments (
  id          BIGSERIAL     NOT NULL,
  content     VARCHAR(255)  NOT NULL CHECK (length(content) > 0),
  author_id   BIGINT        NOT NULL,
  post_id     BIGINT        NOT NULL,
  created_at  TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
  updated_at  TIMESTAMPTZ   NOT NULL DEFAULT NOW(),

  PRIMARY KEY (id),
  FOREIGN KEY (author_id) REFERENCES users (id) ON DELETE CASCADE,
  FOREIGN KEY (post_id)   REFERENCES posts (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS post_likes (
  user_id BIGINT NOT NULL,
  post_id BIGINT NOT NULL,

  PRIMARY KEY (user_id, post_id),
  FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
  FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS comment_likes (
  user_id     BIGINT NOT NULL,
  comment_id  BIGINT NOT NULL,

  PRIMARY KEY (user_id, comment_id),
  FOREIGN KEY (user_id)    REFERENCES users (id) ON DELETE CASCADE,
  FOREIGN KEY (comment_id) REFERENCES comments (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS users_subscriptions (
  user_id         BIGINT NOT NULL,
  subscription_id BIGINT NOT NULL,

  PRIMARY KEY (user_id, subscription_id),
  FOREIGN KEY (user_id)         REFERENCES users (id) ON DELETE CASCADE,
  FOREIGN KEY (subscription_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE VIEW post_info
AS
SELECT
    p.id,
    p.author_id,
    p.content,
    p.created_at,
    p.updated_at,
    COALESCE(lc.count_likes, 0) as likes
FROM posts p
LEFT JOIN (
    SELECT post_id, COUNT(*) as count_likes
    FROM post_likes
    GROUP BY post_id
) as lc ON p.id = lc.post_id;

CREATE VIEW comment_info
AS
SELECT
    c.id,
    c.author_id,
    c.post_id,
    c.content,
    c.created_at,
    c.updated_at,
    COALESCE(lc.count_likes, 0) as likes
FROM comments c
LEFT JOIN (
    SELECT comment_id, COUNT(*) as count_likes
    FROM comment_likes
    GROUP BY comment_id
) lc ON c.id = lc.comment_id;
