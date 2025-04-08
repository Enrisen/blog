-- Posts Table
CREATE TABLE posts (
    post_id SERIAL PRIMARY KEY,            -- Auto-incrementing ID for the post
    author_id INT NOT NULL,                -- Foreign key linking to the user who wrote the post
    title VARCHAR(255) NOT NULL,          -- Title of the blog post
    content TEXT NOT NULL,                -- The main body content of the post
    excerpt TEXT NULL,                     -- A short summary or teaser for the post. MIght have it be AI generated
    view_count INT NOT NULL DEFAULT 0,    -- Counter for how many times the post has been viewed
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- When the post was created/published

    -- Foreign Key constraint: Ensure author_id refers to a valid user
    CONSTRAINT fk_author
        FOREIGN KEY(author_id)
        REFERENCES users(user_id)
        ON DELETE RESTRICT -- Prevent deleting the admin user if they have posts
);