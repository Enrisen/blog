-- Images Table
-- Stores references to images associated with posts.
CREATE TABLE images (
    image_id SERIAL PRIMARY KEY,           -- Auto-incrementing ID for the image
    post_id INT NOT NULL,                -- Foreign key linking to the post this image belongs to
    file_path VARCHAR(512) NOT NULL UNIQUE, -- Path or URL where the image file is stored. Must be unique.

    -- Foreign Key constraint: Ensure post_id refers to a valid post
    CONSTRAINT fk_post_image
        FOREIGN KEY(post_id)
        REFERENCES posts(post_id)
        ON DELETE CASCADE -- If a post is deleted, automatically delete associated image records
);