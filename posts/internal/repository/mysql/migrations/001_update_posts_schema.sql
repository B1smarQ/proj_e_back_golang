-- Update posts table to use UUIDs and add author_name
ALTER TABLE posts
    MODIFY COLUMN id VARCHAR(36) NOT NULL,
    MODIFY COLUMN author_id VARCHAR(36) NOT NULL,
    MODIFY COLUMN thread_id VARCHAR(36) NOT NULL,
    ADD COLUMN author_name VARCHAR(255) NOT NULL AFTER author_id;

-- Drop existing primary key and recreate with UUID
ALTER TABLE posts DROP PRIMARY KEY;
ALTER TABLE posts ADD PRIMARY KEY (id);

-- Add indexes
CREATE INDEX idx_posts_author_id ON posts(author_id);
CREATE INDEX idx_posts_thread_id ON posts(thread_id);
CREATE INDEX idx_posts_creation_time ON posts(creation_time); 