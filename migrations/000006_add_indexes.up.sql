-- Recommended Indexes for performance
CREATE INDEX idx_posts_author_id ON posts(author_id);
CREATE INDEX idx_posts_created_at ON posts(created_at DESC); -- For ordering posts chronologically

CREATE INDEX idx_post_categories_post_id ON post_categories(post_id); -- Useful for finding categories for a post
CREATE INDEX idx_post_categories_category_id ON post_categories(category_id); -- Useful for finding posts in a category

CREATE INDEX idx_images_post_id ON images(post_id);