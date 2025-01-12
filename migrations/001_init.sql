
CREATE TABLE IF NOT EXISTS users (
    id BIGINT PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS segments(
    id SERIAL PRIMARY KEY,
    slug TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS user_segments(
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    segment_id BIGINT NOT NULL REFERENCES segments(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, segment_id)
);

CREATE INDEX idx_slug ON segments(slug);
CREATE INDEX idx_user_id ON user_segments(user_id);
CREATE INDEX idx_segment_id ON user_segments(segment_id);
CREATE INDEX idx_user_segments ON user_segments(user_id, segment_id);


CREATE TABLE IF NOT EXISTS user_segments_history(
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    segment_slug VARCHAR(255) NOT NULL,
    operation_type VARCHAR(50) NOT NULL,
    operation_date TIMESTAMP NOT NULL DEFAULT NOW()
);

ALTER TABLE user_segments ADD COLUMN ttl TIMESTAMP NULL;

