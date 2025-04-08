-- Post Categories
-- Links posts to their categories.
CREATE TABLE post_categories (
    post_id INT NOT NULL,
    category_id INT NOT NULL,

    -- Ensure a post is linked to a category only once
    PRIMARY KEY (post_id, category_id),

    -- Foreign Key constraint: Link to posts table
    CONSTRAINT fk_post
        FOREIGN KEY(post_id)
        REFERENCES posts(post_id)
        ON DELETE CASCADE, -- If a post is deleted, remove its category links

    -- Foreign Key constraint: Link to categories table
    CONSTRAINT fk_category
        FOREIGN KEY(category_id)
        REFERENCES categories(category_id)
        ON DELETE CASCADE -- If a category is deleted, remove links from posts
);